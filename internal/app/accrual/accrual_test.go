package accrual

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"iryzzh/practicum-gophermart/cmd/gophermart/config"
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store/pgstore"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

type reward struct {
	Match      string `json:"match"`
	Reward     int    `json:"reward"`
	RewardType string `json:"reward_type"`
}

type goods struct {
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}

type item struct {
	OrderNumber string  `json:"order"`
	Goods       []goods `json:"goods"`
}

func TestAccrual(t *testing.T) {
	type want struct {
		status string
	}
	tests := []struct {
		name   string
		reward reward
		item   item

		want want
	}{
		{
			name: "ok",
			want: want{
				status: "PROCESSED",
			},
			reward: reward{
				Match:      "Euro",
				Reward:     13,
				RewardType: "%",
			},
			item: item{
				OrderNumber: model.TestOrderNew(t, 1).Number,
				Goods: []goods{
					{
						Description: "Euro",
						Price:       1596.30,
					},
				},
			},
		},
	}

	cfg, err := config.New()
	assert.NoError(t, err)

	var binary string
	switch runtime.GOOS {
	case "windows":
		binary = "accrual_windows_amd64.exe"
	case "linux":
		binary = "accrual_linux_amd64"
	default:
		binary = "accrual_darwin_amd64"
	}

	s, teardown := pgstore.TestDB(t, cfg.DatabaseURI)
	defer teardown("orders")

	client := New(s, cfg.AccrualSystemAddress, time.Second)

	path, _ := os.Getwd()
	binary, _ = filepath.Abs(fmt.Sprintf("%s/../../../cmd/accrual/%s", path, binary))
	_, err = os.Stat(binary)
	assert.NoError(t, err)

	addr, err := url.Parse(cfg.AccrualSystemAddress)
	assert.NoError(t, err)

	cmd := exec.Command(binary, "-a", addr.Host, "-d", cfg.DatabaseURI)

	go func() {
		err := cmd.Run()
		assert.NoError(t, err)
	}()

	time.Sleep(1 * time.Second)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint := fmt.Sprintf("%s/api/goods", cfg.AccrualSystemAddress)
			body, err := json.Marshal(tt.reward)
			assert.NoError(t, err)

			buffer := bytes.NewBuffer(body)
			resp, err := http.Post(endpoint, "application/json", buffer)
			assert.NoError(t, err)
			defer resp.Body.Close()

			endpoint = fmt.Sprintf("%s/api/orders", cfg.AccrualSystemAddress)

			body, err = json.Marshal(tt.item)
			assert.NoError(t, err)
			buffer = bytes.NewBuffer(body)
			resp, err = http.Post(endpoint, "application/json", buffer)
			assert.NoError(t, err)
			defer resp.Body.Close()

			order := model.TestOrderNew(t, 1)
			order.UploadedAt = model.Time{Time: time.Now()}

			assert.NoError(t, s.Order().Create(order))

			err = client.accrualOrderInfo(order)
			assert.NoError(t, err)

			assert.Equal(t, tt.want.status, order.Status.String())

			assert.NoError(t, s.Order().Update(order))
		})
	}
}

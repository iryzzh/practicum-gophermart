package handler_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"iryzzh/practicum-gophermart/internal/app/handler"
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store"
	"iryzzh/practicum-gophermart/internal/app/store/memstore"
	"iryzzh/practicum-gophermart/internal/utils"
	"log"
	"net/http"
	"testing"
)

func TestHandler_GetBalanceHandler(t *testing.T) {
	endpoint := "/api/user/balance"
	withdrawEP := endpoint + "/withdraw"

	st := memstore.New()
	defer st.Close()

	ts, err := handler.NewTestServer(t, st)
	assert.NoError(t, err)
	defer ts.Close()

	u := model.TestUser(t)
	err = st.User().Create(u)
	assert.NoError(t, err)

	err = store.TestOrderWithAccrual(st, u.ID, 5)
	assert.NoError(t, err)

	// login
	jar, err := handler.TestLogin(t, ts, u.Login, u.Password)
	assert.NoError(t, err)

	// проверяем получение текущего баланса
	resp, body := handler.TestRequest(t, "GET", ts.URL+endpoint, nil, jar)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	var m model.Balance
	err = json.Unmarshal([]byte(body), &m)
	assert.NoError(t, err)

	// проверяем списание баланса

	randomSum := func() int {
		res, _ := st.Balance().Get(u.ID)

		min := int(res.Current) / 5
		max := int(res.Current) / 2

		return utils.Intn(max-min+1) + min
	}

	tests := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid withdraw",
			payload: map[string]interface{}{
				"order": utils.RandLuhn(10),
				"sum":   randomSum(),
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid withdraw - not enough funds",
			payload: map[string]interface{}{
				"order": utils.RandLuhn(10),
				"sum":   randomSum() * 10000,
			},
			expectedCode: http.StatusPaymentRequired,
		},
		{
			name: "invalid withdraw - incorrect order number",
			payload: map[string]interface{}{
				"order": "123",
				"sum":   0,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			err = json.NewEncoder(b).Encode(tt.payload)
			assert.NoError(t, err)

			resp, _ = handler.TestRequest(t, "POST", ts.URL+withdrawEP, b, jar)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
			defer func() {
				_ = resp.Body.Close()
			}()
		})
	}

	// проверяем withdrawals

	resp, body = handler.TestRequest(t, "GET", ts.URL+"/api/user/withdrawals", nil, jar)
	log.Println("response body:", body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()
}

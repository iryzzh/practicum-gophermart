package handler_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"iryzzh/practicum-gophermart/internal/app/handler"
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store"
	"iryzzh/practicum-gophermart/internal/app/store/memstore"
	"iryzzh/practicum-gophermart/internal/utils"
	"net/http"
	"strings"
	"testing"
)

func TestHandler_GetOrdersHandler(t *testing.T) {
	tests := []struct {
		name               string
		user               *model.User
		orders             []*model.Order
		expectedStatusCode int
	}{
		{
			name: "get order ok",
			user: model.TestUser(t),
			orders: []*model.Order{
				{
					Number:     utils.RandLuhn(10),
					Status:     model.OrderNew,
					UploadedAt: model.TestRandomDate(t),
				},
				{
					Number:     utils.RandLuhn(10),
					Status:     model.OrderNew,
					UploadedAt: model.TestRandomDate(t),
				},
				{
					Number:     utils.RandLuhn(10),
					Status:     model.OrderNew,
					UploadedAt: model.TestRandomDate(t),
				},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "get order no content",
			user: &model.User{
				Login:    "test-user2",
				Password: "password2",
			},
			orders:             nil,
			expectedStatusCode: http.StatusNoContent,
		},
	}

	st := memstore.New()
	defer st.Close()

	ts, err := handler.NewTestServer(t, st)
	assert.NoError(t, err)
	defer ts.Close()

	url := ts.URL + "/api/user/orders"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := st.User().Create(tt.user); err != nil {
				assert.Equal(t, store.ErrUserAlreadyExists, err)
			}

			// присваиваем существующий айди пользователя к заказам
			// и создаем эти заказы в базе данных
			for _, order := range tt.orders {
				order.UserID = tt.user.ID
				assert.NoError(t, st.Order().Create(order))
			}

			jar, err := handler.TestLogin(t, ts, tt.user.Login, tt.user.Password)
			assert.NoError(t, err)

			resp, body := handler.TestRequest(t, "GET", url, nil, jar)
			defer resp.Body.Close()
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			// заказов нет
			if tt.orders == nil {
				return
			}

			var orders []*model.Order
			err = json.Unmarshal([]byte(body), &orders)
			assert.NoError(t, err)

			for i, o := range orders {
				// проверка сортировки дат
				if i > 0 {
					if o.UploadedAt.After(orders[i-1].UploadedAt.Time) || o.UploadedAt.After(orders[0].UploadedAt.Time) {
						assert.Fail(t, "incorrect date sorting")
					}
				}
			}
		})
	}
}

func TestHandler_PostOrdersHandler(t *testing.T) {
	endpoint := "/api/user/orders"

	st := memstore.New()
	defer st.Close()

	// создаем пользователя
	u := model.TestUser(t)
	if err := st.User().Create(u); err != nil {
		t.Fatal(err)
	}

	ts, err := handler.NewTestServer(t, st)
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()

	// получаем куки
	jar, err := handler.TestLogin(t, ts, u.Login, u.Password)
	assert.NoError(t, err)

	tests := []struct {
		name         string
		payload      string
		expectedCode int
	}{
		{
			name:         "valid order",
			payload:      "12345678903",
			expectedCode: http.StatusAccepted,
		},
		{
			name:         "valid order 2",
			payload:      "12345678903",
			expectedCode: http.StatusOK,
		},
		{
			name:         "conflict order",
			payload:      "12345678903",
			expectedCode: http.StatusConflict,
		},
		{
			name:         "invalid order",
			payload:      "123456789",
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			directOrderCreate := func(orderID string, userID int) error {
				if exists, _ := st.Order().FindByNumber(orderID); exists == nil {
					if err := st.Order().Create(&model.Order{
						Number: orderID,
						UserID: userID,
					}); err != nil {
						return err
					}
				}

				return nil
			}

			switch tt.expectedCode {
			case http.StatusOK:
				err = directOrderCreate(tt.payload, u.ID)
				assert.NoError(t, err)
			case http.StatusConflict:
				err = directOrderCreate(tt.payload, u.ID)
				assert.NoError(t, err)

				u2 := &model.User{Login: "random-login", Password: "random-password"}
				err = st.User().Create(u2)
				assert.NoError(t, err)

				jar, err = handler.TestLogin(t, ts, u2.Login, u2.Password)
				assert.NoError(t, err)
			}

			// создаем заказ
			resp, _ := handler.TestRequest(t, "POST", ts.URL+endpoint, strings.NewReader(tt.payload), jar)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
			_ = resp.Body.Close()
		})
	}
}

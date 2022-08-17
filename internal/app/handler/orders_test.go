package handler_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"iryzzh/practicum-gophermart/internal/app/handler"
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store"
	"iryzzh/practicum-gophermart/internal/app/store/memstore"
	"net/http"
	"strings"
	"testing"
)

func TestHandler_GetOrdersHandler(t *testing.T) {
	endpoint := "/api/user/orders"

	st := memstore.New()
	defer st.Close()

	ts, err := handler.NewTestServer(t, st)
	assert.NoError(t, err)
	defer ts.Close()

	// создаем пользователя
	u := model.TestUser(t)
	err = st.User().Create(u)
	assert.NoError(t, err)

	// получаем куки
	jar, err := handler.TestLogin(t, ts, u.Login, u.Password)
	assert.NoError(t, err)

	// проверяем пустой ответ
	resp, _ := handler.TestRequest(t, "GET", ts.URL+endpoint, nil, jar)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	_ = resp.Body.Close()

	// создаем тестовые заказы
	err = store.TestOrder(st, u, 5)
	assert.NoError(t, err)

	// проверяем заказ через веб запрос
	resp, body := handler.TestRequest(t, "GET", ts.URL+endpoint, nil, jar)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	var orders []model.Order
	err = json.Unmarshal([]byte(body), &orders)
	assert.NoError(t, err)

	for i, v := range orders {
		if v.Accrual != nil && *v.Accrual == float32(0) {
			assert.Fail(t, "accrual should not be zero")
		}

		if v.UploadedAt.String() == "" {
			assert.Fail(t, "empty time")
		}

		if i > 0 {
			if orders[i].UploadedAt.After(orders[i-1].UploadedAt.Time) {
				assert.Fail(t, "incorrect sorting")
			}
		}
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

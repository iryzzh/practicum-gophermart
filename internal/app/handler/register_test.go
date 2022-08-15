package handler_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"iryzzh/practicum-gophermart/internal/app/handler"
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store/memstore"
	"net/http"
	"net/http/cookiejar"
	"testing"
)

func TestHandler_RegisterHandler(t *testing.T) {
	endpoint := "/api/user/register"
	st := memstore.New()
	ts, err := handler.NewTestServer(t, st)
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()

	tests := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"login":    "login",
				"password": "password",
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "too short password",
			payload: map[string]interface{}{
				"login":    "random_login",
				"password": "1",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "conflict user name",
			payload:      model.TestUser(t),
			expectedCode: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			err := json.NewEncoder(b).Encode(tt.payload)
			if err != nil {
				t.Fatal(err)
			}

			if tt.expectedCode == http.StatusConflict {
				if err := st.User().Create(model.TestUser(t)); err != nil {
					t.Fatal(err)
				}
			}

			jar, err := cookiejar.New(nil)
			if err != nil {
				t.Fatal(err)
			}

			resp, _ := handler.TestRequest(t, "POST", ts.URL+endpoint, b, jar)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			if resp.StatusCode == http.StatusOK {
				assert.NotNil(t, jar)
			}
		})
	}
}

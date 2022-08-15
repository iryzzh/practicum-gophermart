package handler_test

import (
	"github.com/stretchr/testify/assert"
	"iryzzh/practicum-gophermart/internal/app/handler"
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store/memstore"
	"testing"
)

func TestHandler_LoginHandler(t *testing.T) {
	st := memstore.New()
	defer st.Close()

	ts, err := handler.NewTestServer(t, st)
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()

	u := model.TestUser(t)
	if err := st.User().Create(u); err != nil {
		t.Fatal(err)
	}

	_, err = handler.TestLogin(t, ts, u.Login, u.Password)
	assert.NoError(t, err)
}

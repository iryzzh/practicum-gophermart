package memstore_test

import (
	"github.com/stretchr/testify/assert"
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store"
	"iryzzh/practicum-gophermart/internal/app/store/memstore"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	st := memstore.New()
	user := model.TestUser(t)

	assert.NoError(t, st.User().Create(user))

	err := st.User().Create(user)
	assert.Equal(t, store.ErrUserAlreadyExists, err)
}

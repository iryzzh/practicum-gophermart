package pgstore_test

import (
	"github.com/stretchr/testify/assert"
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store"
	"iryzzh/practicum-gophermart/internal/app/store/pgstore"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	s, teardown := pgstore.TestDB(t, dsn)
	defer teardown("users")
	//_ = teardown

	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))

	err := s.User().Create(u)
	assert.Equal(t, store.ErrUserAlreadyExists, err)
	u2, err := s.User().FindByLogin(u.Login)
	assert.NoError(t, err)
	assert.Equal(t, u.ID, u2.ID)
}

package pgstore_test

import (
	"github.com/stretchr/testify/assert"
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store"
	"iryzzh/practicum-gophermart/internal/app/store/pgstore"
	"testing"
)

func TestOrderRepository(t *testing.T) {
	s, teardown := pgstore.TestDB(t, dsn)
	defer teardown("orders, users")

	order1 := model.TestOrderNew(t, 1)
	order2 := model.TestOrderProcessed(t, 1)
	assert.NoError(t, s.Order().Create(order1))
	assert.NoError(t, s.Order().Update(order1))
	assert.NoError(t, s.Order().Create(order2))
	assert.NoError(t, s.Order().Update(order2))

	assert.Equal(t, store.ErrOrderAlreadyExists, s.Order().Create(order1))
	assert.Equal(t, store.ErrOrderAlreadyExists, s.Order().Create(order2))

	r, err := s.Order().FindByNumber(order1.Number)
	assert.NoError(t, err)
	assert.Equal(t, order1.Number, r.Number)
	assert.Equal(t, order1.ID, r.ID)
	assert.Equal(t, order1.UserID, r.UserID)
	assert.Equal(t, order1.Status, r.Status)
	assert.Equal(t, order1.Accrual, r.Accrual)
	//assert.Equal(t, order1.UploadedAt.String(), r.UploadedAt.String())

	r, err = s.Order().FindByNumber(order2.Number)
	assert.NoError(t, err)
	assert.Equal(t, order2.Number, r.Number)
	assert.Equal(t, order2.ID, r.ID)
	assert.Equal(t, order2.UserID, r.UserID)
	assert.Equal(t, order2.Status, r.Status)
	assert.Equal(t, order2.Accrual, r.Accrual)
	//assert.Equal(t, order2.UploadedAt.String(), r.UploadedAt.String())

	r, err = s.Order().FindByID(order1.ID)
	assert.NoError(t, err)
	assert.Equal(t, order1.Number, r.Number)
	assert.Equal(t, order1.ID, r.ID)
	assert.Equal(t, order1.UserID, r.UserID)
	assert.Equal(t, order1.Status, r.Status)
	assert.Equal(t, order1.Accrual, r.Accrual)
}

package pgstore_test

import (
	"github.com/stretchr/testify/assert"
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store/pgstore"
	"testing"
)

func TestBalanceRepository_Get(t *testing.T) {
	s, teardown := pgstore.TestDB(t, dsn)
	defer teardown("orders")

	order := model.TestOrderProcessed(t, 1)
	assert.NoError(t, s.Order().Create(order))
	assert.NoError(t, s.Order().Update(order))

	balance, err := s.Balance().Get(1)
	assert.NoError(t, err)
	assert.Equal(t, *order.Accrual, balance.Current)
}

func TestBalanceRepository_Withdraw(t *testing.T) {
	s, teardown := pgstore.TestDB(t, dsn)
	defer teardown("orders, withdrawals")

	order := model.TestOrderProcessed(t, 1)
	assert.NoError(t, s.Order().Create(order))
	assert.NoError(t, s.Order().Update(order))

	withdraw := &model.Withdraw{
		UserID:      1,
		Sum:         *order.Accrual,
		OrderNumber: order.Number,
	}

	assert.NoError(t, s.Balance().Withdraw(withdraw))
	assert.Error(t, s.Balance().Withdraw(withdraw))

	balance, err := s.Balance().Get(1)
	assert.NoError(t, err)

	data, err := s.Balance().Withdrawals(1)
	assert.NoError(t, err)
	assert.Equal(t, data[0].Sum, balance.Withdrawn)
}

package memstore_test

import (
	"github.com/stretchr/testify/assert"
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store/memstore"
	"testing"
)

func TestBalanceRepository_Get(t *testing.T) {
	st := memstore.New()
	defer st.Close()

	order := model.TestOrderProcessed(t, 1)
	assert.NoError(t, st.Order().Create(order))

	balance, err := st.Balance().Get(1)
	assert.NoError(t, err)
	assert.Equal(t, balance.Current, *order.Accrual)
}

func TestBalanceRepository_Withdraw(t *testing.T) {
	st := memstore.New()
	defer st.Close()

	order := model.TestOrderProcessed(t, 1)
	assert.NoError(t, st.Order().Create(order))

	withdraw := &model.Withdraw{
		UserID:      1,
		Sum:         *order.Accrual,
		OrderNumber: order.Number,
	}

	assert.NoError(t, st.Balance().Withdraw(withdraw))
	assert.Error(t, st.Balance().Withdraw(withdraw))

	balance, err := st.Balance().Get(1)
	assert.NoError(t, err)

	data, err := st.Balance().Withdrawals(1)
	assert.NoError(t, err)
	assert.Equal(t, data[0].Sum, balance.Withdrawn)
}

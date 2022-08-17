package memstore_test

import (
	"github.com/stretchr/testify/assert"
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store"
	"iryzzh/practicum-gophermart/internal/app/store/memstore"
	"testing"
)

func TestOrderRepository(t *testing.T) {
	st := memstore.New()
	defer st.Close()

	order1 := model.TestOrderNew(t, 1)
	order2 := model.TestOrderProcessed(t, 1)
	assert.NoError(t, st.Order().Create(order1))
	assert.NoError(t, st.Order().Create(order2))

	assert.Equal(t, store.ErrOrderAlreadyExists, st.Order().Create(order1))
	assert.Equal(t, store.ErrOrderAlreadyExists, st.Order().Create(order2))

	r, err := st.Order().FindByNumber(order1.Number)
	assert.NoError(t, err)
	assert.Equal(t, r, order1)

	r, err = st.Order().FindByNumber(order2.Number)
	assert.NoError(t, err)
	assert.Equal(t, r, order2)

	r, err = st.Order().FindByID(order1.ID)
	assert.NoError(t, err)
	assert.Equal(t, r, order1)

	orders, err := st.Order().GetByUserID(1)
	assert.NoError(t, err)

	var f1, f2 bool
	for _, v := range orders {
		if v.Number == order1.Number {
			if assert.Equal(t, v, order1) {
				f1 = true
			}
		}
		if v.Number == order2.Number {
			if assert.Equal(t, v, order2) {
				f2 = true
			}
		}
	}
	assert.Condition(t, func() bool {
		if f1 && f2 {
			return true
		}

		return false
	})
}

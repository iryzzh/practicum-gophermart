package store

import (
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/utils"
	"time"
)

func TestOrderWithAccrual(store Store, userID int, count int) error {
	if count == 0 {
		count = 1
	}

	for i := 0; i < count; i++ {
		min := 100
		max := 500

		accrual := float32(utils.Intn(max-min+1) + min)

		n := utils.RandLuhn(10)

		order := &model.Order{
			UserID:     userID,
			Number:     n,
			Accrual:    &accrual,
			Status:     model.OrderProcessed,
			UploadedAt: model.Time{Time: time.Now().Add(-24 * time.Duration(utils.Intn(max-min+1)+min) * time.Hour).Add(-time.Duration(utils.Intn(max-min+1)+min) * 27 * time.Minute)},
		}

		err := store.Order().Create(order)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestOrder(store Store, u *model.User, count int) error {
	if count == 0 {
		count = 1
	}

	for i := 0; i < count; i++ {
		min := 10
		max := 30
		var accrual *float32
		var status model.OrderStatus
		if utils.Intn(max-min+1)+min < 20 {
			accrual = func(f float32) *float32 {
				return &f
			}(float32((utils.Intn(max-min+1) + min) * 25))

			status = model.OrderProcessed
		}

		n := utils.RandLuhn(10)
		order := &model.Order{
			UserID:     u.ID,
			Number:     n,
			Accrual:    accrual,
			Status:     status,
			UploadedAt: model.Time{Time: time.Now().Add(-24 * time.Duration(utils.Intn(max-min+1)+min) * time.Hour).Add(-time.Duration(utils.Intn(max-min+1)+min) * 27 * time.Minute)},
		}

		err := store.Order().Create(order)
		if err != nil {
			return err
		}
	}

	return nil
}

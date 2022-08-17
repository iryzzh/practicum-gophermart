package memstore

import (
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store"
	"sort"
	"time"
)

type BalanceRepository struct {
	store *Store
}

func (b *BalanceRepository) Withdrawals(userID int) ([]*model.Withdraw, error) {
	var withdraws []*model.Withdraw

	b.store.Lock()
	defer b.store.Unlock()

	for _, v := range b.store.withdraws {
		if v.UserID == userID {
			withdraws = append(withdraws, v)
		}
	}

	if len(withdraws) == 0 {
		return nil, store.ErrRecordNotFound
	}

	sort.Slice(withdraws, func(i, j int) bool {
		return withdraws[i].ProcessedAt.String() > withdraws[j].ProcessedAt.String()
	})

	return withdraws, nil
}

func (b *BalanceRepository) Withdraw(w *model.Withdraw) error {
	balance, err := b.Get(w.UserID)
	if err != nil {
		return err
	}

	if balance.Current < w.Sum {
		return store.ErrNotEnoughFunds
	}

	b.store.Lock()
	defer b.store.Unlock()

	w.ID = b.store.withdrawNextID + 1

	// testing:
	if w.ProcessedAt.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
		w.ProcessedAt = model.Time{Time: time.Now()}
	}

	b.store.withdraws[b.store.withdrawNextID] = w
	b.store.withdrawNextID++

	return nil
}

func (b *BalanceRepository) Get(userID int) (*model.Balance, error) {
	b.store.Lock()
	defer b.store.Unlock()

	var accrualTotal float32
	for _, v := range b.store.orders {
		if v.UserID == userID {
			if *v.Accrual > float32(0) {
				accrualTotal += *v.Accrual
			}
		}
	}

	var withdrawsTotal float32
	for _, v := range b.store.withdraws {
		if v.UserID == userID {
			if v.Sum > float32(0) {
				withdrawsTotal += v.Sum
			}
		}
	}

	return &model.Balance{
		Current:   accrualTotal - withdrawsTotal,
		Withdrawn: withdrawsTotal,
	}, nil
}

package pgstore

import (
	"database/sql"
	"errors"
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store"
	"sort"
)

type BalanceRepository struct {
	store *Store
}

func (b *BalanceRepository) Get(userID int) (*model.Balance, error) {
	var accrualTotal float32
	err := b.store.db.QueryRow("select COALESCE( sum(accrual), 0 ) from orders where user_id = $1", userID).Scan(&accrualTotal)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	var withdrawalsTotal float32
	err = b.store.db.QueryRow("select COALESCE( sum(withdraw), 0 ) from withdrawals where user_id = $1", userID).Scan(&withdrawalsTotal)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	return &model.Balance{
		Current:   accrualTotal - withdrawalsTotal,
		Withdrawn: withdrawalsTotal,
	}, nil
}

func (b *BalanceRepository) Withdraw(w *model.Withdraw) error {
	balance, err := b.Get(w.UserID)
	if err != nil {
		return err
	}

	if balance.Current < w.Sum {
		return store.ErrNotEnoughFunds
	}

	return b.store.db.QueryRow(
		"insert into withdrawals (user_id, order_number, withdraw, processed_at) values ($1, $2, $3, $4) returning id",
		w.UserID, w.OrderNumber, w.Sum, w.ProcessedAt.String()).Scan(&w.ID)
}

func (b *BalanceRepository) Withdrawals(userID int) ([]*model.Withdraw, error) {
	var withdrawals []*model.Withdraw

	rows, err := b.store.db.Query(
		"select id, user_id, order_number, withdraw, processed_at from withdrawals where user_id = $1",
		userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for rows.Next() {
		var withdraw model.Withdraw
		err := rows.Scan(&withdraw.ID, &withdraw.UserID, &withdraw.OrderNumber, &withdraw.Sum, &withdraw.ProcessedAt)
		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, &withdraw)
	}

	sort.Slice(withdrawals, func(i, j int) bool {
		return withdrawals[i].ProcessedAt.String() > withdrawals[j].ProcessedAt.String()
	})

	return withdrawals, nil
}

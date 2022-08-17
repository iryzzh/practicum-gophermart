package memstore

import (
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store"
	"sync"
)

type Store struct {
	sync.RWMutex

	users          map[int]*model.User
	orders         map[int]*model.Order
	withdraws      map[int]*model.Withdraw
	userNextID     int
	orderNextID    int
	withdrawNextID int
}

func (s *Store) Close() error {
	return nil
}

func New() *Store {
	return &Store{
		users:          make(map[int]*model.User),
		orders:         make(map[int]*model.Order),
		withdraws:      make(map[int]*model.Withdraw),
		userNextID:     0,
		orderNextID:    0,
		withdrawNextID: 0,
	}
}

func (s *Store) User() store.UserRepository {
	return &UserRepository{store: s}
}

func (s *Store) Order() store.OrderRepository {
	return &OrderRepository{store: s}
}

func (s *Store) Balance() store.BalanceRepository {
	return &BalanceRepository{store: s}
}

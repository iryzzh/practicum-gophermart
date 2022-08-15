package store

import "iryzzh/practicum-gophermart/internal/app/model"

type Store interface {
	User() UserRepository
	Order() OrderRepository
	Balance() BalanceRepository
	Close() error
}

type UserRepository interface {
	Create(user *model.User) error
	FindByID(id int) (*model.User, error)
	FindByLogin(login string) (*model.User, error)
}

type OrderRepository interface {
	Exists(number string, userID int) error
	Create(order *model.Order) error
	Update(order *model.Order) error
	FindByID(id int) (*model.Order, error)
	FindByUserID(userID int) (*model.Order, error)
	FindByNumber(number string) (*model.Order, error)
	GetByUserID(userID int) ([]*model.Order, error)
	Incomplete() ([]*model.Order, error)
}

type BalanceRepository interface {
	Get(userID int) (*model.Balance, error)
	Withdraw(w *model.Withdraw) error
	Withdrawals(userID int) ([]*model.Withdraw, error)
}

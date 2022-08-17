package store

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrRecordNotFound     = errors.New("record not found")
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrOrderConflict      = errors.New("order conflict")
	ErrOrderUpdate        = errors.New("order update error")
	ErrNotEnoughFunds     = errors.New("not enough funds in the account")
)

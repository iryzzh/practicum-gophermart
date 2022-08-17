package memstore

import (
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store"
)

type UserRepository struct {
	store *Store
}

func (u *UserRepository) Create(user *model.User) error {
	if err := user.Validate(0); err != nil {
		return err
	}

	if err := user.BeforeCreate(); err != nil {
		return err
	}

	if u, _ := u.store.User().FindByLogin(user.Login); u != nil {
		return store.ErrUserAlreadyExists
	}

	u.store.Lock()
	defer u.store.Unlock()

	user.ID = u.store.userNextID + 1

	u.store.users[u.store.userNextID] = user
	u.store.userNextID++

	return nil
}

func (u *UserRepository) FindByLogin(login string) (*model.User, error) {
	u.store.Lock()
	defer u.store.Unlock()

	for _, v := range u.store.users {
		if login == v.Login && !v.Deleted {
			return v, nil
		}
	}

	return nil, store.ErrUserNotFound
}

func (u *UserRepository) FindByID(id int) (*model.User, error) {
	u.store.Lock()
	defer u.store.Unlock()

	for _, v := range u.store.users {
		if id == v.ID && !v.Deleted {
			return v, nil
		}
	}

	return nil, store.ErrUserNotFound
}

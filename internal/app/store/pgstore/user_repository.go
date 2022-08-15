package pgstore

import (
	"database/sql"
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

	if u, _ := u.FindByLogin(user.Login); u != nil {
		return store.ErrUserAlreadyExists
	}

	return u.store.db.QueryRow(
		"INSERT INTO users (login, encrypted_password) values ($1, $2) returning id",
		user.Login, user.EncryptedPassword).Scan(&user.ID)
}

func (u *UserRepository) FindByID(id int) (*model.User, error) {
	user := &model.User{}

	if err := u.store.db.QueryRow(
		"SELECT id, login, encrypted_password, deleted, deleted_at FROM users WHERE id = $1",
		id,
	).Scan(
		&user.ID,
		&user.Login,
		&user.EncryptedPassword,
		&user.Deleted,
		&user.DeletedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return user, nil
}

func (u *UserRepository) FindByLogin(login string) (*model.User, error) {
	user := &model.User{}

	if err := u.store.db.QueryRow(
		"SELECT id, login, encrypted_password, deleted, deleted_at from users WHERE login = $1",
		login,
	).Scan(
		&user.ID,
		&user.Login,
		&user.EncryptedPassword,
		&user.Deleted,
		&user.DeletedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return user, nil
}

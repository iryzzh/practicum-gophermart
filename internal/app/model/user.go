package model

import (
	"database/sql"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

var DefaultMinPasswordLength = 6

type User struct {
	ID                int          `json:"id"`
	Login             string       `json:"login"`
	Password          string       `json:"password,omitempty"`
	EncryptedPassword string       `json:"-"`
	Deleted           bool         `json:"_"`
	DeletedAt         sql.NullTime `json:"-"`
}

func (u *User) Validate(minPasswordLength int) error {
	if minPasswordLength <= 0 {
		minPasswordLength = DefaultMinPasswordLength
	}

	return validation.ValidateStruct(
		u,
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")),
			validation.Length(minPasswordLength, 100)),
	)
}

// BeforeCreate шифрует поле пароля и сохраняет его в поле "EncryptedPassword"
func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	}

	return nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// Sanitize очищает поле пароля
func (u *User) Sanitize() {
	u.Password = ""
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

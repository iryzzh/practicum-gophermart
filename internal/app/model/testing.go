package model

import (
	"iryzzh/practicum-gophermart/internal/utils"
	"testing"
	"time"
)

func TestUser(t *testing.T) *User {
	t.Helper()

	return &User{
		Login:    "test-user",
		Password: "password",
	}
}

func TestOrderNew(t *testing.T, userID int) *Order {
	t.Helper()

	return &Order{
		UserID: userID,
		Number: "12345678903",
		Status: OrderNew,
	}
}

func TestRandomDate(t *testing.T) Time {
	min := 100
	max := 500

	return Time{Time: time.Now().Add(-24 * time.Duration(utils.Intn(max-min+1)+min) * time.Hour).Add(-time.Duration(utils.Intn(max-min+1)+min) * 27 * time.Minute)}
}

func TestOrderProcessed(t *testing.T, userID int) *Order {
	t.Helper()

	accrual := float32(500)

	return &Order{
		UserID:  userID,
		Number:  "9278923470",
		Status:  OrderProcessed,
		Accrual: &accrual,
	}
}

package model

import "testing"

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

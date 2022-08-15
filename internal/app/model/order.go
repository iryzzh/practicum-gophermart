package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"strings"
)

type OrderStatus int

const (
	OrderNew OrderStatus = iota
	OrderProcessing
	OrderInvalid
	OrderProcessed
)

func (s *OrderStatus) String() string {
	switch *s {
	case OrderNew:
		return "NEW"
	case OrderProcessing:
		return "PROCESSING"
	case OrderInvalid:
		return "INVALID"
	case OrderProcessed:
		return "PROCESSED"
	default:
		return "UNKNOWN"
	}
}

func (s *OrderStatus) UnmarshalOrderStatus(str string) OrderStatus {
	switch str {
	case "NEW":
		return OrderNew
	case "PROCESSING":
		return OrderProcessing
	case "INVALID":
		return OrderInvalid
	case "PROCESSED":
		return OrderProcessed
	default:
		return 0
	}
}

func (s *OrderStatus) Value() (driver.Value, error) {
	return fmt.Sprintf("%v", *s), nil
}

func (s *OrderStatus) Scan(value interface{}) error {
	if v, ok := value.([]uint8); ok {
		*s = s.UnmarshalOrderStatus(string(v))
		return nil
	}

	return fmt.Errorf("can't convert %T to OrderStatus", value)
}

func (s *OrderStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *OrderStatus) UnmarshalJSON(b []byte) error {
	data := strings.ReplaceAll(string(b), "\"", "")
	switch data {
	case "PROCESSING":
		*s = OrderProcessing
	case "INVALID":
		*s = OrderInvalid
	case "PROCESSED":
		*s = OrderProcessed
	default:
		*s = OrderNew
	}
	return nil
}

type Order struct {
	ID         int          `json:"-"`
	UserID     int          `json:"-"`
	Number     string       `json:"number"`
	Status     OrderStatus  `json:"status"`
	Accrual    *float32     `json:"accrual,omitempty"`
	UploadedAt Time         `json:"uploaded_at"`
	Deleted    bool         `json:"-"`
	DeletedAt  sql.NullTime `json:"-"`
}

func (order *Order) Validate() error {
	return validation.ValidateStruct(
		order,
		validation.Field(&order.Number, validation.Required, validation.By(isValidLuhn)),
	)
}

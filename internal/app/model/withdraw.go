package model

import validation "github.com/go-ozzo/ozzo-validation"

type Withdraw struct {
	ID          int     `json:"-"`
	UserID      int     `json:"-"`
	OrderNumber string  `json:"order"`
	Sum         float32 `json:"sum"`
	ProcessedAt Time    `json:"processed_at"`
}

func (w *Withdraw) Validate() error {
	return validation.ValidateStruct(
		w,
		validation.Field(&w.OrderNumber, validation.By(isValidLuhn)),
	)
}

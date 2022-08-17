package model

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"strconv"
)

func requiredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}

		return nil
	}
}

func isValidLuhn(value interface{}) error {
	number, err := strconv.Atoi(value.(string))
	if err != nil {
		return err
	}
	if (number%10+checksumLuhn(number/10))%10 != 0 {
		return errors.New("invalid order number")
	}

	return nil
}

func checksumLuhn(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}

	return luhn % 10
}

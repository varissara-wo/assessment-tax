package admin

import "errors"

const (
	ErrInvalidPersonalGreaterAmount = "amount must be greater than 10000.0"
	ErrInvalidPersonalLessAmount    = "amount must be less than 100000.0"
)

func (a Amount) ValidatePersonalDeduction() error {
	if a.Amount < 10000 {
		return errors.New(ErrInvalidPersonalGreaterAmount)
	}
	if a.Amount > 100000 {
		return errors.New(ErrInvalidPersonalLessAmount)
	}
	return nil
}

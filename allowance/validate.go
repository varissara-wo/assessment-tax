package allowance

import "errors"

const (
	ErrInvalidKReceiptLessAmount    = "amount must be greater than 0.0"
	ErrInvalidKReceiptGreaterAmount = "amount must be less than 100000.0"
)

func (a Amount) ValidateKReceipt() error {
	if a.Amount <= 0 {
		return errors.New(ErrInvalidKReceiptLessAmount)
	}

	if a.Amount > 100000 {
		return errors.New(ErrInvalidKReceiptGreaterAmount)
	}

	return nil
}

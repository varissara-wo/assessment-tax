package allowance

import "errors"

const (
	ErrInvalidPersonalGreaterAmount = "amount must be greater than 10000.0"
	ErrInvalidPersonalLessAmount    = "amount must be less than 100000.0"
	ErrInvalidKReceiptGreaterAmount = "amount must be greater than 0.0"
	ErrInvalidKReceiptLessAmount    = "amount must be less than 100000.0"
	ErrInvalidAllowance             = "allowances must be donation and k-receipt only"
	ErrInvalidAllowanceAmount       = "allowance amount must be greater than or equal to 0"
)

func (a Amount) ValidatePersonal() error {
	if a.Amount < 10000 {
		return errors.New(ErrInvalidPersonalGreaterAmount)
	}
	if a.Amount > 100000 {
		return errors.New(ErrInvalidPersonalLessAmount)
	}
	return nil
}

func (a Amount) ValidateKReceipt() error {
	if a.Amount <= 0 {
		return errors.New(ErrInvalidKReceiptGreaterAmount)
	}

	if a.Amount > 100000 {
		return errors.New(ErrInvalidKReceiptLessAmount)
	}

	return nil
}

func ValidateAllowance(a Allowance) error {
	if err := validateAllowanceType(a); err != nil {
		return err
	}

	if a.Amount < 0 {
		return errors.New(ErrInvalidAllowanceAmount)
	}

	return nil
}

func validateAllowanceType(a Allowance) error {
	for _, validType := range ValidAllowanceTypes {
		if a.AllowanceType == validType {
			return nil
		}
	}
	return errors.New(ErrInvalidAllowance)
}

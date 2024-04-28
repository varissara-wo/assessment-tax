package tax

import (
	"errors"

	"github.com/varissara-wo/assessment-tax/allowance"
)

func (td *TaxDetails) ValidateTaxDetails() error {
	if err := validateTotalIncome(td.TotalIncome); err != nil {
		return err
	}

	if err := validateWHT(td.WHT, td.TotalIncome); err != nil {
		return err
	}

	for _, a := range td.Allowances {
		if err := validateAllowance(a); err != nil {
			return err
		}
	}

	return nil
}

func validateTotalIncome(i float64) error {
	if i < 0 {
		return errors.New(ErrInvalidTotalIncome)
	}
	return nil
}

func validateWHT(wht, totalIncome float64) error {
	if wht < 0 || wht > totalIncome {
		return errors.New(ErrInvalidWHT)
	}
	return nil
}

func validateAllowance(a allowance.Allowance) error {
	if err := validateAllowanceType(a); err != nil {
		return err
	}

	if a.Amount < 0 {
		return errors.New(ErrInvalidAllowanceAmount)
	}

	return nil
}

func validateAllowanceType(a allowance.Allowance) error {
	for _, validType := range allowance.ValidAllowanceTypes {
		if a.AllowanceType == validType {
			return nil
		}
	}
	return errors.New(ErrInvalidAllowance)
}

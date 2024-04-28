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
		if err := allowance.ValidateAllowance(a); err != nil {
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

package tax

import "errors"

func (td *TaxDetails) ValidateTaxDetails() error {
	if err := validateTotalIncome(td.TotalIncome); err != nil {
		return err
	}

	if err := validateWHT(td.Wht, td.TotalIncome); err != nil {
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

func validateAllowance(a Allowance) error {
	if err := validateAllowanceType(a); err != nil {
		return err
	}

	if a.Amount < 0 {
		return errors.New(ErrInvalidAllowanceAmount)
	}

	return nil
}

func validateAllowanceType(a Allowance) error {
	for _, validType := range validAllowanceTypes {
		if a.AllowanceType == validType {
			return nil
		}
	}
	return errors.New(ErrInvalidAllowance)
}

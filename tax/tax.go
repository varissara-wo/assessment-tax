package tax

import "errors"

type Allowance struct {
	AllowanceType string
	Amount        float64
}

type TaxDetails struct {
	TotalIncome float64
	Wht         float64
	Allowances  []Allowance
}

type Tax struct {
	Tax string `json:"tax"`
}

const (
	ErrInvalidTotalIncome = "total income must be greater than 0"
	ErrInvalidWht         = "wht must be greater than or equal to 0 and less than total income"
	ErrInvalidAllowance   = "allowances must be donation and k-receipt only"
	ErrInvalidAmount      = "allowance amount must be greater than or equal to 0"
)

func (td *TaxDetails) ValidateTaxDetails() error {
	if td.TotalIncome <= 0 {
		return errors.New(ErrInvalidTotalIncome)
	}

	if td.Wht < 0 || td.Wht > td.TotalIncome {
		return errors.New(ErrInvalidWht)
	}

	for _, a := range td.Allowances {
		if a.AllowanceType != "donation" && a.AllowanceType != "k-receiptt" {
			return errors.New(ErrInvalidAllowance)
		}

		if a.Amount < 0 {
			return errors.New(ErrInvalidAmount)
		}
	}

	return nil
}

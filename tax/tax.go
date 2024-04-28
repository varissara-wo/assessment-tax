package tax

import "github.com/varissara-wo/assessment-tax/allowance"

type TaxDetails struct {
	TotalIncome float64
	WHT         float64
	Allowances  []allowance.Allowance
}

type TaxBreakdown struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type TaxCalculationResponse struct {
	Tax      float64        `json:"tax"`
	TaxLevel []TaxBreakdown `json:"taxLevel"`
}

const (
	ErrInvalidTotalIncome     = "total income must be greater than or equals 0"
	ErrInvalidWHT             = "wht must be greater than or equal to 0 and less than total income"
	ErrInvalidAllowance       = "allowances must be donation and k-receipt only"
	ErrInvalidAllowanceAmount = "allowance amount must be greater than or equal to 0"
)

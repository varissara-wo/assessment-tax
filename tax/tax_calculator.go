package tax

import (
	"math"

	"github.com/varissara-wo/assessment-tax/allowance"
)

type TaxBracket struct {
	Description string
	MaxIncome   float64
	TaxRate     float64
	MaxTax      float64
}

var taxBrackets = []TaxBracket{
	{Description: "0-150,000", MaxIncome: 150000.0, TaxRate: 0.0, MaxTax: 0.0},
	{Description: "150,001-500,000", MaxIncome: 500000.0, TaxRate: 0.1, MaxTax: 35000.0},
	{Description: "500,001-1,000,000", MaxIncome: 1000000.0, TaxRate: 0.15, MaxTax: 75000.0},
	{Description: "1,000,001-2,000,000", MaxIncome: 2000000.0, TaxRate: 0.2, MaxTax: 200000.0},
	{Description: "2,000,001 ขึ้นไป", MaxIncome: math.MaxFloat64, TaxRate: 0.35},
}

func CalculateTax(income float64, wht float64) TaxResponse {
	tax := 0.0
	previousMaxTax := 0.0
	previousMaxIncome := 0.0
	tbl := []TaxBreakdown{}

	for _, bracket := range taxBrackets {
		var tb TaxBreakdown
		if income <= bracket.MaxIncome && income > previousMaxIncome {
			tax = ((income - previousMaxIncome) * bracket.TaxRate) + previousMaxTax
			tb = TaxBreakdown{
				Level: bracket.Description,
				Tax:   (income - previousMaxIncome) * bracket.TaxRate,
			}
		} else {
			if income > bracket.MaxIncome {
				tb = TaxBreakdown{
					Level: bracket.Description,
					Tax:   bracket.MaxTax,
				}
			} else {
				tb = TaxBreakdown{
					Level: bracket.Description,
					Tax:   0.0,
				}
			}

		}
		tbl = append(tbl, tb)
		previousMaxTax += bracket.MaxTax
		previousMaxIncome = bracket.MaxIncome
	}

	if tax-wht < 0 {
		r := TaxResponse{
			TaxRefund: wht - tax,
			TaxLevel:  tbl,
		}
		return r
	}

	r := TaxResponse{
		Tax:      tax - wht,
		TaxLevel: tbl,
	}

	return r
}

func (td TaxDetails) CalculateNetIncome(ma allowance.MaxAllowance) float64 {
	return td.TotalIncome - allowance.CalculateAllowances(td.Allowances, ma)
}

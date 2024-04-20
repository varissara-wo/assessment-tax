package tax

import (
	"math"
)

type TaxBracket struct {
	MaxIncome float64
	TaxRate   float64
	MaxTax    float64
}

var taxBrackets = []TaxBracket{
	{MaxIncome: 150000.0, TaxRate: 0.0, MaxTax: 0.0},
	{MaxIncome: 500000.0, TaxRate: 0.1, MaxTax: 35000.0},
	{MaxIncome: 1000000.0, TaxRate: 0.15, MaxTax: 110000.0},
	{MaxIncome: 2000000.0, TaxRate: 0.2, MaxTax: 310000.0},
	{MaxIncome: math.MaxFloat64, TaxRate: 0.35},
}

func CalculateTax(income float64) float64 {
	var tax float64
	var previousMaxTax float64
	var previousMaxIncome float64

	for _, bracket := range taxBrackets {
		if income <= bracket.MaxIncome {
			tax = ((income - previousMaxIncome) * bracket.TaxRate) + previousMaxTax
			break
		}
		previousMaxTax = bracket.MaxTax
		previousMaxIncome = bracket.MaxIncome
	}
	return math.Round(tax)
}

type AllowanceAmount struct {
	Donation float64
	KReceipt float64
	Personal float64
}

type MaxAllowance struct {
	Donation float64
	KReceipt float64
	Personal float64
}

func calculateAllowances(allowances []Allowance, ma MaxAllowance) float64 {
	aa := AllowanceAmount{
		Donation: 0.0,
		KReceipt: 0.0,
		Personal: 0.0,
	}

	aa.Personal = ma.Personal

	for _, a := range allowances {
		switch a.AllowanceType {
		case Donation:
			if aa.Donation+a.Amount <= ma.Donation {
				aa.Donation += a.Amount
			} else {
				aa.Donation = ma.Donation
			}
		case KReceipt:
			if aa.KReceipt+a.Amount <= ma.KReceipt {
				aa.KReceipt += a.Amount
			} else {
				aa.KReceipt = ma.KReceipt
			}
		}
	}

	return aa.Donation + aa.KReceipt + aa.Personal
}

func (td TaxDetails) CalculateNetIncome(ma MaxAllowance) float64 {
	return td.TotalIncome - calculateAllowances(td.Allowances, ma) - td.Wht
}

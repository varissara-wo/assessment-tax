package postgres

import (
	"fmt"
	"math"

	"github.com/varissara-wo/assessment-tax/tax"
)

const (
	MAX_INCOME_1 = 150000.0
	TAX_RATE_1   = 0.0

	MAX_INCOME_2 = 500000.0
	TAX_RATE_2   = 0.1
	MAX_TAX_2    = 35000.0

	MAX_INCOME_3 = 1000000.0
	TAX_RATE_3   = 0.15
	MAX_TAX_3    = 110000.0

	MAX_INCOME_4 = 2000000.0
	TAX_RATE_4   = 0.2
	MAX_TAX_4    = 310000.0

	TAX_RATE_5 = 0.35
)

func calculateTax(income float64) float64 {
	if income <= MAX_INCOME_1 {
		return income * TAX_RATE_1
	} else if income > MAX_INCOME_1 && income <= MAX_INCOME_2 {
		return math.Round((income - MAX_INCOME_1) * TAX_RATE_2)
	} else if income > MAX_INCOME_2 && income <= MAX_INCOME_3 {
		return math.Round(MAX_TAX_2 + (income-MAX_INCOME_2)*TAX_RATE_3)
	} else if income > MAX_INCOME_3 && income <= MAX_INCOME_4 {
		return math.Round(MAX_TAX_3 + (income-MAX_INCOME_3)*TAX_RATE_4)
	}
	return math.Round(MAX_TAX_4 + (income-MAX_INCOME_4)*TAX_RATE_5)
}

func calculateAllowances(allowances []tax.Allowance) float64 {
	personal := 60000.0
	donation := 0.0
	kreceip := 0.0

	for _, a := range allowances {
		if a.AllowanceType == "donation" {
			donation += a.Amount
		} else if a.AllowanceType == "k-receipt" {
			kreceip += a.Amount
		}
	}

	return personal + donation + kreceip
}

func calculateNetIncome(td tax.TaxDetails) float64 {
	return td.TotalIncome - calculateAllowances(td.Allowances) - td.Wht
}

func (p *Postgres) TaxCalculation(td tax.TaxDetails) (tax.Tax, error) {

	netIncome := calculateNetIncome(td)

	if netIncome <= 0 {
		return tax.Tax{Tax: fmt.Sprintf("%.1f", 0.0)}, nil
	}

	taxAmount := calculateTax(calculateNetIncome(td))

	return tax.Tax{Tax: fmt.Sprintf("%.1f", taxAmount)}, nil
}

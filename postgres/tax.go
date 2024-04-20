package postgres

import (
	"fmt"
	"math"

	"github.com/varissara-wo/assessment-tax/tax"
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

func calculateTax(income float64) float64 {
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

	taxAmount := calculateTax(netIncome)

	return tax.Tax{Tax: fmt.Sprintf("%.1f", taxAmount)}, nil
}

package postgres

import (
	"fmt"

	"github.com/varissara-wo/assessment-tax/tax"
)

func (p *Postgres) TaxCalculation(td tax.TaxDetails) (tax.Tax, error) {

	netIncome := td.CalculateNetIncome()

	if netIncome <= 0 {
		return tax.Tax{Tax: fmt.Sprintf("%.1f", 0.0)}, nil
	}

	taxAmount := tax.CalculateTax(netIncome)

	return tax.Tax{Tax: fmt.Sprintf("%.1f", taxAmount)}, nil
}

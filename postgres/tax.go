package postgres

import (
	"github.com/varissara-wo/assessment-tax/tax"
)

func (p *Postgres) TaxCalculation(td tax.TaxDetails) (tax.TaxResponse, error) {

	ma, err := p.GetAllowances()
	if err != nil {
		return tax.TaxResponse{}, err
	}

	netIncome := td.CalculateNetIncome(ma)

	return tax.CalculateTax(netIncome, td.WHT), nil
}

func (p *Postgres) TaxesCalculation(tds []tax.TaxDetails) ([]tax.Taxes, error) {

	ma, err := p.GetAllowances()
	if err != nil {
		return []tax.Taxes{}, err
	}

	taxes := []tax.Taxes{}

	for _, td := range tds {

		if err := td.ValidateTaxDetails(); err != nil {
			return []tax.Taxes{}, err
		}

		netIncome := td.CalculateNetIncome(ma)

		result := tax.CalculateTax(netIncome, td.WHT)

		tr := tax.Taxes{
			TotalIncome: td.TotalIncome,
			Tax:         result.Tax,
			TaxRefund:   result.TaxRefund,
		}

		taxes = append(taxes, tr)

	}

	return taxes, nil
}

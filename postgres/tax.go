package postgres

import (
	"database/sql"

	"github.com/varissara-wo/assessment-tax/allowance"
	"github.com/varissara-wo/assessment-tax/tax"
)

func (p *Postgres) TaxCalculation(td tax.TaxDetails) (tax.TaxCalculationResponse, error) {

	var rows *sql.Rows

	rows, err := p.Db.Query("SELECT * FROM allowances")

	if err != nil {
		return tax.TaxCalculationResponse{}, err
	}

	defer rows.Close()

	var ma allowance.MaxAllowance

	for rows.Next() {
		var t allowance.AllowanceType
		var amount float64
		var id int
		err := rows.Scan(&id, &t, &amount)

		if err != nil {
			return tax.TaxCalculationResponse{}, err
		}

		switch t {
		case allowance.Donation:
			ma.Donation = amount
		case allowance.KReceipt:
			ma.KReceipt = amount
		case allowance.Personal:
			ma.Personal = amount
		}
	}

	netIncome := td.CalculateNetIncome(ma)

	return tax.CalculateTax(netIncome, td.WHT), nil
}

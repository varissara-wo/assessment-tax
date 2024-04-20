package postgres

import (
	"database/sql"
	"fmt"

	"github.com/varissara-wo/assessment-tax/tax"
)

func (p *Postgres) TaxCalculation(td tax.TaxDetails) (tax.TaxCalculationResponse, error) {

	var rows *sql.Rows

	rows, err := p.Db.Query("SELECT * FROM allowances")

	if err != nil {
		return tax.TaxCalculationResponse{}, err
	}

	defer rows.Close()

	var ma tax.MaxAllowance

	for rows.Next() {
		var t tax.AllowanceType
		var amount float64
		var id int
		err := rows.Scan(&id, &t, &amount)
		if err != nil {
			return tax.TaxCalculationResponse{}, err
		}

		switch t {
		case tax.Donation:
			ma.Donation = amount
		case tax.KReceipt:
			ma.KReceipt = amount
		case tax.Personal:
			ma.Personal = amount
		}
	}

	netIncome := td.CalculateNetIncome(ma)

	if netIncome <= 0 {
		return tax.TaxCalculationResponse{Tax: fmt.Sprintf("%.1f", 0.0)}, nil
	}

	taxAmount := tax.CalculateTax(netIncome)

	return tax.TaxCalculationResponse{Tax: fmt.Sprintf("%.1f", taxAmount), TaxLevel: []tax.TaxBreakdown{}}, nil
}

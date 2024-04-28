package postgres

import (
	"github.com/varissara-wo/assessment-tax/allowance"
)

func (p *Postgres) GetAllowances() (allowance.MaxAllowance, error) {
	var ma allowance.MaxAllowance

	rows, err := p.Db.Query("SELECT * FROM allowances")
	if err != nil {
		return ma, err
	}
	defer rows.Close()

	for rows.Next() {
		var t allowance.AllowanceType
		var amount float64
		var id int
		err := rows.Scan(&id, &t, &amount)
		if err != nil {
			return ma, err
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

	return ma, nil
}

func (p *Postgres) SetPersonal(a float64) (allowance.PersonalDeduction, error) {
	_, err := p.Db.Exec("UPDATE allowances SET max_amount = $1 WHERE type = 'personal'", a)
	if err != nil {
		return allowance.PersonalDeduction{}, err
	}
	return allowance.PersonalDeduction{Personal: a}, nil
}

func (p *Postgres) SetKReceipt(a float64) (allowance.KReceiptDeduction, error) {
	_, err := p.Db.Exec("UPDATE allowances SET max_amount = $1 WHERE type = 'k-receipt'", a)
	if err != nil {
		return allowance.KReceiptDeduction{}, err
	}
	return allowance.KReceiptDeduction{KReceipt: a}, nil
}

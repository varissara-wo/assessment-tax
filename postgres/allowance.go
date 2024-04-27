package postgres

import (
	"github.com/varissara-wo/assessment-tax/admin"
	"github.com/varissara-wo/assessment-tax/allowance"
	"github.com/varissara-wo/assessment-tax/tax"
)

func (p *Postgres) GetAllowances() (tax.MaxAllowance, error) {
	var ma tax.MaxAllowance

	rows, err := p.Db.Query("SELECT * FROM allowances")
	if err != nil {
		return ma, err
	}
	defer rows.Close()

	for rows.Next() {
		var t tax.AllowanceType
		var amount float64
		var id int
		err := rows.Scan(&id, &t, &amount)
		if err != nil {
			return ma, err
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

	return ma, nil
}

func (p *Postgres) SetPersonalAllowance(a float64) (admin.PersonalAllowance, error) {
	_, err := p.Db.Exec("UPDATE allowances SET max_amount = $1 WHERE type = 'personal'", a)
	if err != nil {
		return admin.PersonalAllowance{}, err
	}
	return admin.PersonalAllowance{PersonalDeduction: a}, nil
}

func (p *Postgres) SetPersonal(a float64) (allowance.Personal, error) {
	_, err := p.Db.Exec("UPDATE allowances SET max_amount = $1 WHERE type = 'personal'", a)
	if err != nil {
		return allowance.Personal{}, err
	}
	return allowance.Personal{Personal: a}, nil
}

func (p *Postgres) SetKReceipt(a float64) (allowance.KReceipt, error) {
	_, err := p.Db.Exec("UPDATE allowances SET max_amount = $1 WHERE type = 'k-receipt'", a)
	if err != nil {
		return allowance.KReceipt{}, err
	}
	return allowance.KReceipt{KReceipt: a}, nil
}

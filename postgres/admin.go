package postgres

import "github.com/varissara-wo/assessment-tax/admin"

func (p *Postgres) SetPersonalDeduction(a float64) (admin.PersonalDeduction, error) {
	_, err := p.Db.Exec("UPDATE allowances SET max_amount = $1 WHERE type = 'personal'", a)
	if err != nil {
		return admin.PersonalDeduction{}, err
	}
	return admin.PersonalDeduction{PersonalDeduction: a}, nil
}

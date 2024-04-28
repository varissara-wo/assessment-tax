package allowance

func CalculateAllowances(allowances []Allowance, ma MaxAllowance) float64 {
	aa := AllowanceAmount{
		Donation: 0.0,
		KReceipt: 0.0,
		Personal: 0.0,
	}

	aa.Personal = ma.Personal

	for _, a := range allowances {
		switch a.AllowanceType {
		case "donation":
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

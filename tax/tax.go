package tax

type AllowanceType string

const (
	Donation  AllowanceType = "donation"
	KReceiptt AllowanceType = "k-receiptt"
)

var validAllowanceTypes = []AllowanceType{Donation, KReceiptt}

type Allowance struct {
	AllowanceType AllowanceType
	Amount        float64
}

type TaxDetails struct {
	TotalIncome float64
	Wht         float64
	Allowances  []Allowance
}

type Tax struct {
	Tax string `json:"tax"`
}

const (
	ErrInvalidTotalIncome     = "total income must be greater than or equals 0"
	ErrInvalidWHT             = "wht must be greater than or equal to 0 and less than total income"
	ErrInvalidAllowance       = "allowances must be donation and k-receipt only"
	ErrInvalidAllowanceAmount = "allowance amount must be greater than or equal to 0"
)

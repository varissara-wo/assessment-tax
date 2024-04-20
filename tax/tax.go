package tax

type AllowanceType string

const (
	Donation AllowanceType = "donation"
	KReceipt AllowanceType = "k-receipt"
	Personal AllowanceType = "personal"
)

var validAllowanceTypes = []AllowanceType{Donation, KReceipt}

type Allowance struct {
	AllowanceType AllowanceType
	Amount        float64
}

type TaxDetails struct {
	TotalIncome float64
	Wht         float64
	Allowances  []Allowance
}

type TaxBreakdown struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type TaxCalculationResponse struct {
	Tax      float64        `json:"tax"`
	TaxLevel []TaxBreakdown `json:"taxLevel"`
}

const (
	ErrInvalidTotalIncome     = "total income must be greater than or equals 0"
	ErrInvalidWHT             = "wht must be greater than or equal to 0 and less than total income"
	ErrInvalidAllowance       = "allowances must be donation and k-receipt only"
	ErrInvalidAllowanceAmount = "allowance amount must be greater than or equal to 0"
)

package allowance

type KReceiptDeduction struct {
	KReceipt float64 `json:"kReceipt"`
}

type PersonalDeduction struct {
	Personal float64 `json:"personalDeduction"`
}

type Amount struct {
	Amount float64 `json:"amount"`
}

type AllowanceType string

const (
	Donation AllowanceType = "donation"
	KReceipt AllowanceType = "k-receipt"
	Personal AllowanceType = "personal"
)

var ValidAllowanceTypes = []AllowanceType{Donation, KReceipt}

type Allowance struct {
	AllowanceType AllowanceType
	Amount        float64
}

type AllowanceAmount struct {
	Donation float64
	KReceipt float64
	Personal float64
}

type MaxAllowance struct {
	Donation float64
	KReceipt float64
	Personal float64
}

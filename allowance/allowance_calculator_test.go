package allowance

import "testing"

var mockMaxAllowance = MaxAllowance{
	Donation: 100000.0,
	KReceipt: 50000.0,
	Personal: 60000.0,
}

func TestAllowancesCalculation(t *testing.T) {
	t.Run("allowance should return 180000.0", func(t *testing.T) {
		want := 180000.0

		mockAllowances := []Allowance{
			{
				AllowanceType: KReceipt,
				Amount:        20000.0,
			},
			{
				AllowanceType: Donation,
				Amount:        105000.0,
			},
		}

		got := CalculateAllowances(mockAllowances, mockMaxAllowance)

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("should not exceed max allowances", func(t *testing.T) {
		allowances := []Allowance{
			{
				AllowanceType: Donation,
				Amount:        100.0,
			},
			{
				AllowanceType: Donation,
				Amount:        200000.0,
			},
			{
				AllowanceType: KReceipt,
				Amount:        100000.0,
			},
		}

		expected := 210000.0

		got := CalculateAllowances(allowances, mockMaxAllowance)

		if got != expected {
			t.Errorf("expected %v but got %v", expected, got)
		}
	})
}

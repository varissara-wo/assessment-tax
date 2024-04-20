package tax

import (
	"fmt"
	"testing"
)

func TestTaxCalculation(t *testing.T) {

	tests := []struct {
		income float64
		tax    float64
	}{
		{0.0, 0.0},
		{150000.0, 0.0},
		{150001.0, 0.0},
		{500000.0, 35000.0},
		{500001.0, 35000.0},
		{1000000.0, 110000.0},
		{1000001.0, 110000.0},
		{2000000.0, 310000.0},
		{3000000.0, 660000.0},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("final income %v should return %v", tt.income, tt.tax), func(t *testing.T) {
			want := tt.tax

			got := CalculateTax(tt.income)

			if got != want {
				t.Errorf("got %v want %v", got, want)
			}
		})
	}
}

var mockMaxAllowance = MaxAllowance{
	Donation: 100000.0,
	KReceipt: 50000.0,
	Personal: 60000.0,
}

func TestAllowancesCalculation(t *testing.T) {
	t.Run("should return 60000", func(t *testing.T) {
		want := 180000.0

		mockAllowances := []Allowance{
			{
				AllowanceType: "k-receipt",
				Amount:        20000.0,
			},
			{
				AllowanceType: "donation",
				Amount:        105000.0,
			},
		}

		got := calculateAllowances(mockAllowances, mockMaxAllowance)

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestNetIncomeCalculation(t *testing.T) {
	t.Run("should return 60000", func(t *testing.T) {
		want := 788000.0

		mockTaxDetails := TaxDetails{
			TotalIncome: 1000000.0,
			Wht:         2000.0,
			Allowances: []Allowance{
				{
					AllowanceType: "k-receipt",
					Amount:        200000.0,
				},
				{
					AllowanceType: "donation",
					Amount:        100000.0,
				},
			},
		}

		got := mockTaxDetails.CalculateNetIncome(mockMaxAllowance)

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

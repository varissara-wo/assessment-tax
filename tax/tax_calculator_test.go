package tax

import (
	"fmt"
	"reflect"
	"testing"
)

func generateTaxBreakdown(taxValues ...float64) []TaxBreakdown {
	var breakdown []TaxBreakdown
	for i, tax := range taxValues {
		breakdown = append(breakdown, TaxBreakdown{
			Level: taxBrackets[i].Description,
			Tax:   tax,
		})
	}
	return breakdown
}

func TestTaxCalculation(t *testing.T) {

	tests := []struct {
		income   float64
		tax      float64
		taxLevel []TaxBreakdown
	}{
		{0.0, 0.0, generateTaxBreakdown(0.0, 0.0, 0.0, 0.0, 0.0)},
		{150000.0, 0.0, generateTaxBreakdown(0.0, 0.0, 0.0, 0.0, 0.0)},
		{150001.0, 0.1, generateTaxBreakdown(0.0, 0.1, 0.0, 0.0, 0.0)},
		{500000.0, 35000.0, generateTaxBreakdown(0.0, 35000.0, 0.0, 0.0, 0.0)},
		{500001.0, 35000.15, generateTaxBreakdown(0.0, 35000.0, 0.15, 0.0, 0.0)},
		{1000000.0, 110000.0, generateTaxBreakdown(0.0, 35000.0, 75000.0, 0.0, 0.0)},
		{1000001.0, 110000.2, generateTaxBreakdown(0.0, 35000.0, 75000.0, 0.2, 0.0)},
		{2000000.0, 310000.0, generateTaxBreakdown(0.0, 35000.0, 75000.0, 200000.0, 0.0)},
		{3000000.0, 660000.0, generateTaxBreakdown(0.0, 35000.0, 75000.0, 200000.0, 350000.0)},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("final income %v should return %v", tt.income, tt.tax), func(t *testing.T) {
			want := tt.tax
			wantTaxLevel := tt.taxLevel

			got := CalculateTax(tt.income, 0.0)

			if got.Tax != want {
				t.Errorf("got %v want %v", got.Tax, want)
			}

			if !reflect.DeepEqual(got.TaxLevel, wantTaxLevel) {
				t.Errorf("got tax level %v want %v", got.TaxLevel, wantTaxLevel)
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
		want := 790000.0

		mockTaxDetails := TaxDetails{
			TotalIncome: 1000000.0,
			WHT:         2000.0,
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

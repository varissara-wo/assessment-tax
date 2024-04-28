package tax

import (
	"errors"
	"testing"

	"github.com/varissara-wo/assessment-tax/allowance"
)

func TestValidateTaxDetails(t *testing.T) {

	t.Run("should return nil if all fields are valid", func(t *testing.T) {
		mockTaxDetails := TaxDetails{
			TotalIncome: 500000.0,
			WHT:         0.0,
			Allowances: []allowance.Allowance{
				{
					AllowanceType: "donation",
					Amount:        0.0,
				},
			},
		}

		err := mockTaxDetails.ValidateTaxDetails()

		if err != nil {
			t.Errorf("expected nil but got %v", err)
		}
	})

	testCases := []struct {
		name          string
		taxDetails    TaxDetails
		expectedError error
	}{

		{
			name: "should return an error if total income is less than 0",
			taxDetails: TaxDetails{
				TotalIncome: -1.0,
				WHT:         0.0,
				Allowances: []allowance.Allowance{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			expectedError: errors.New(ErrInvalidTotalIncome),
		},
		{
			name: "should return an error if WHT is less than 0",
			taxDetails: TaxDetails{
				TotalIncome: 500000.0,
				WHT:         -1.0,
				Allowances: []allowance.Allowance{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			expectedError: errors.New(ErrInvalidWHT),
		},
		{
			name: "should return an error if WHT is greater than total income",
			taxDetails: TaxDetails{
				TotalIncome: 500000.0,
				WHT:         500001.0,
				Allowances: []allowance.Allowance{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			expectedError: errors.New(ErrInvalidWHT),
		},
		{
			name: "should return an error if allowance type is not donation or k-receipt",
			taxDetails: TaxDetails{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []allowance.Allowance{
					{
						AllowanceType: "invalid",
						Amount:        0.0,
					},
				},
			},
			expectedError: errors.New(allowance.ErrInvalidAllowance),
		},
		{
			name: "should return an error if allowance amount is less than 0",
			taxDetails: TaxDetails{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []allowance.Allowance{
					{
						AllowanceType: "donation",
						Amount:        -1.0,
					},
				},
			},
			expectedError: errors.New(allowance.ErrInvalidAllowanceAmount),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.taxDetails.ValidateTaxDetails()

			if err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error %v but got %v", tc.expectedError, err)
			}
		})
	}

}

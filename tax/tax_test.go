package tax

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestValidateTaxDetails(t *testing.T) {

	t.Run("should return nil if all fields are valid", func(t *testing.T) {
		mockTaxDetails := TaxDetails{
			TotalIncome: 500000.0,
			Wht:         0.0,
			Allowances: []Allowance{
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
				Wht:         0.0,
				Allowances: []Allowance{
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
				Wht:         -1.0,
				Allowances: []Allowance{
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
				Wht:         500001.0,
				Allowances: []Allowance{
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
				Wht:         0.0,
				Allowances: []Allowance{
					{
						AllowanceType: "invalid",
						Amount:        0.0,
					},
				},
			},
			expectedError: errors.New(ErrInvalidAllowance),
		},
		{
			name: "should return an error if allowance amount is less than 0",
			taxDetails: TaxDetails{
				TotalIncome: 500000.0,
				Wht:         0.0,
				Allowances: []Allowance{
					{
						AllowanceType: "donation",
						Amount:        -1.0,
					},
				},
			},
			expectedError: errors.New(ErrInvalidAllowanceAmount),
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

func TestAllowancesCalculation(t *testing.T) {
	t.Run("should return 60000", func(t *testing.T) {
		want := 360000.0

		mockAllowances := []Allowance{
			{
				AllowanceType: "k-receipt",
				Amount:        200000.0,
			},
			{
				AllowanceType: "donation",
				Amount:        100000.0,
			},
		}

		got := calculateAllowances(mockAllowances)

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestNetIncomeCalculation(t *testing.T) {
	t.Run("should return 60000", func(t *testing.T) {
		want := 638000.0

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

		got := mockTaxDetails.CalculateNetIncome()

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

type stub struct {
	TaxDetails TaxDetails
	Tax        Tax
	err        error
}

func (s *stub) TaxCalculation(td TaxDetails) (Tax, error) {
	return s.Tax, s.err
}

func TestTaxHandler(t *testing.T) {
	t.Run("should return 400 and an error if provide bad request payload", func(t *testing.T) {

		mockTaxDetails := TaxDetails{
			TotalIncome: -1.0,
			Wht:         0.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        0.0,
				},
			},
		}

		mockTaxDetailsJSON, _ := json.Marshal(mockTaxDetails)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(mockTaxDetailsJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		st := stub{}

		p := New(&st)
		err := p.TaxHandler(c)

		if err != nil {
			t.Errorf("expected error message but got %v", err)
		}

		var gotErr Err
		json.Unmarshal(rec.Body.Bytes(), &gotErr)

		if gotErr.Message != ErrInvalidTotalIncome {
			t.Errorf("expected error message %v but got %v", ErrInvalidTotalIncome, gotErr.Message)
		}

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status code %v but got %v", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("should return 500 and an error message if the tax calculation fails", func(t *testing.T) {
		mockTaxDetails := TaxDetails{
			TotalIncome: 10000.0,
			Wht:         0.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        0.0,
				},
			},
		}

		mockTaxDetailsJSON, _ := json.Marshal(mockTaxDetails)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(mockTaxDetailsJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		st := stub{err: errors.New("tax calculation fails")}
		p := New(&st)
		err := p.TaxHandler(c)

		if err != nil {
			t.Errorf("got some error %v", err)
		}

		var gotErr Err
		json.Unmarshal(rec.Body.Bytes(), &gotErr)

		if gotErr.Message != st.err.Error() {
			t.Errorf("expected error message %v but got %v", st.err, err)
		}

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %v but got %v", http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("should return 200 and a tax of 29000.0 if the income is 500000.0", func(t *testing.T) {
		mockTaxDetails := TaxDetails{
			TotalIncome: 500000.0,
			Wht:         0.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        0.0,
				},
			},
		}

		e := echo.New()
		mockTaxDetailsJSON, _ := json.Marshal(mockTaxDetails)
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(mockTaxDetailsJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		want := Tax{Tax: "29000.0"}

		st := stub{
			TaxDetails: mockTaxDetails,
			Tax:        want,
		}

		p := New(&st)
		err := p.TaxHandler(c)

		if err != nil {
			t.Errorf("got some error %v", err)
		}

		var got Tax
		json.Unmarshal(rec.Body.Bytes(), &got)

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

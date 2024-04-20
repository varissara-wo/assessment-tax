package tax

import (
	"errors"
	"testing"
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

type stubTax struct {
	TaxDetails TaxDetails
	Tax        Tax
	err        error
}

func (s *stubTax) TaxCalculation(td TaxDetails) (Tax, error) {
	return s.Tax, s.err
}

// func TestTaxHandler(t *testing.T) {
// 	t.Run("should return 400 and an error if total income less than or equals 0", func(t *testing.T) {
// 		mockTaxDetails := TaxDetails{
// 			TotalIncome: -1.0,
// 			Wht:         0.0,
// 			Allowances: []Allowance{
// 				{
// 					AllowanceType: "donation",
// 					Amount:        0.0,
// 				},
// 			},
// 		}

// 		mockTaxDetailsJSON, _ := json.Marshal(mockTaxDetails)

// 		e := echo.New()
// 		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(mockTaxDetailsJSON))
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)

// 		st := stubTax{
// 			TaxDetails: mockTaxDetails,
// 			err:        errors.New(ErrInvalidTotalIncome),
// 		}
// 		p := New(&st)
// 		err := p.TaxHandler(c)

// 		if err != nil {
// 			t.Errorf("got some error %v", err)
// 		}

// 		var gotErr Err
// 		json.Unmarshal(rec.Body.Bytes(), &gotErr)

// 		if gotErr.Message != st.err.Error() {
// 			t.Errorf("expected error message %v but got %v", st.err, err)
// 		}

// 		if rec.Code != http.StatusBadRequest {
// 			t.Errorf("expected status code %v but got %v", http.StatusBadRequest, rec.Code)
// 		}

// 	})

// 	t.Run("should return 500 and an error message if the tax calculation fails", func(t *testing.T) {
// 		e := echo.New()
// 		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBufferString(`{"totalIncome": 5000}`))
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)

// 		st := stubTax{err: errors.New("tax calculation fails")}
// 		p := New(&st)
// 		err := p.TaxHandler(c)

// 		if err != nil {
// 			t.Errorf("got some error %v", err)
// 		}

// 		var gotErr Err
// 		json.Unmarshal(rec.Body.Bytes(), &gotErr)

// 		if gotErr.Message != st.err.Error() {
// 			t.Errorf("expected error message %v but got %v", st.err, err)
// 		}

// 		if rec.Code != http.StatusInternalServerError {
// 			t.Errorf("expected status code %v but got %v", http.StatusInternalServerError, rec.Code)
// 		}
// 	})

// 	t.Run("should return 400 and an error if WHT less than 0 or less than total income", func(t *testing.T) {
// 		mockTaxDetails := TaxDetails{
// 			TotalIncome: 500000.0,
// 			Wht:         500001.0,
// 			Allowances: []Allowance{
// 				{
// 					AllowanceType: "donation",
// 					Amount:        0.0,
// 				},
// 			},
// 		}

// 		mockTaxDetailsJSON, _ := json.Marshal(mockTaxDetails)

// 		e := echo.New()
// 		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(mockTaxDetailsJSON))
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)

// 		st := stubTax{
// 			TaxDetails: mockTaxDetails,
// 			err:        errors.New(ErrInvalidWht),
// 		}
// 		p := New(&st)
// 		err := p.TaxHandler(c)

// 		if err != nil {
// 			t.Errorf("got some error %v", err)
// 		}

// 		var gotErr Err
// 		json.Unmarshal(rec.Body.Bytes(), &gotErr)

// 		if gotErr.Message != st.err.Error() {
// 			t.Errorf("expected error message %v but got %v", st.err, err)
// 		}

// 		if rec.Code != http.StatusBadRequest {
// 			t.Errorf("expected status code %v but got %v", http.StatusBadRequest, rec.Code)
// 		}
// 	})

// 	t.Run("should return a tax of 29000.0 if the income is 500000.0", func(t *testing.T) {
// 		mockTaxDetails := TaxDetails{
// 			TotalIncome: 500000.0,
// 			Wht:         0.0,
// 			Allowances: []Allowance{
// 				{
// 					AllowanceType: "donation",
// 					Amount:        0.0,
// 				},
// 			},
// 		}

// 		e := echo.New()
// 		mockTaxDetailsJSON, _ := json.Marshal(mockTaxDetails)
// 		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(mockTaxDetailsJSON))
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)

// 		want := Tax{Tax: "29000.0"}

// 		st := stubTax{
// 			TaxDetails: mockTaxDetails,
// 			Tax:        want,
// 		}

// 		p := New(&st)
// 		err := p.TaxHandler(c)

// 		if err != nil {
// 			t.Errorf("got some error %v", err)
// 		}

// 		var got Tax
// 		json.Unmarshal(rec.Body.Bytes(), &got)

// 		if got != want {
// 			t.Errorf("got %v want %v", got, want)
// 		}
// 	})

// }

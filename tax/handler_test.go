package tax

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/varissara-wo/assessment-tax/allowance"
)

type stub struct {
	TaxDetails TaxDetails
	Tax        TaxResponse
	err        error
}

func (s *stub) TaxCalculation(td TaxDetails) (TaxResponse, error) {
	return s.Tax, s.err
}

func (s *stub) TaxesCalculation(tds []TaxDetails) ([]Taxes, error) {
	return []Taxes{}, nil
}

func TestTaxHandler(t *testing.T) {
	t.Run("should return 400 and an error if provide bad request payload", func(t *testing.T) {

		mockTaxDetails := TaxDetails{
			TotalIncome: -1.0,
			WHT:         0.0,
			Allowances: []allowance.Allowance{
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
			WHT:         0.0,
			Allowances: []allowance.Allowance{
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
			WHT:         0.0,
			Allowances: []allowance.Allowance{
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

		want := TaxResponse{Tax: 29000.0, TaxLevel: []TaxBreakdown{}}

		st := stub{
			TaxDetails: mockTaxDetails,
			Tax:        want,
		}

		p := New(&st)
		err := p.TaxHandler(c)

		if err != nil {
			t.Errorf("got some error %v", err)
		}

		var got TaxResponse
		json.Unmarshal(rec.Body.Bytes(), &got)

		if got.Tax != want.Tax || !reflect.DeepEqual(got.TaxLevel, want.TaxLevel) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

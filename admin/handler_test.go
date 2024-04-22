package admin

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
)

type stub struct {
	Amount            Amount
	PersonalDeduction PersonalDeduction
	err               error
}

func (s *stub) SetPersonalDeduction(amount float64) (PersonalDeduction, error) {
	return s.PersonalDeduction, s.err
}

func TestTaxHandler(t *testing.T) {
	t.Run("should return 400 and error message if amount is less than 10000", func(t *testing.T) {
		mockAmount := Amount{Amount: 9999.0}
		mockAmountJSON, _ := json.Marshal(mockAmount)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/personal-deduction", bytes.NewBuffer(mockAmountJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		st := stub{}

		p := New(&st)
		err := p.PersonalDeductionHandler(c)

		if err != nil {
			t.Errorf("expected error message but got %v", err)
		}

		var gotErr Err
		json.Unmarshal(rec.Body.Bytes(), &gotErr)

		if gotErr.Message != ErrInvalidPersonalGreaterAmount {
			t.Errorf("expected error message %v but got %v", ErrInvalidPersonalGreaterAmount, gotErr.Message)
		}

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status code %v but got %v", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("should return 200 and personal deduction if amount is valid", func(t *testing.T) {
		mockAmount := Amount{Amount: 20000.0}
		mockAmountJSON, _ := json.Marshal(mockAmount)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/personal-deduction", bytes.NewBuffer(mockAmountJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		st := stub{
			PersonalDeduction: PersonalDeduction{
				PersonalDeduction: 20000.0,
			},
		}

		p := New(&st)
		err := p.PersonalDeductionHandler(c)

		if err != nil {
			t.Errorf("expected nil but got %v", err)
		}

		var got PersonalDeduction
		json.Unmarshal(rec.Body.Bytes(), &got)

		want := PersonalDeduction{
			PersonalDeduction: 20000.0,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("expected status code %v but got %v", http.StatusOK, rec.Code)
		}
	})

	t.Run("should return 500 and error message if can't update personal deduction", func(t *testing.T) {
		mockAmount := Amount{Amount: 20000.0}
		mockAmountJSON, _ := json.Marshal(mockAmount)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/personal-deduction", bytes.NewBuffer(mockAmountJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		st := stub{
			err: errors.New("failed to update personal deduction"),
		}

		p := New(&st)
		err := p.PersonalDeductionHandler(c)

		if err != nil {
			t.Errorf("expected nil but got %v", err)
		}

		var gotErr Err
		json.Unmarshal(rec.Body.Bytes(), &gotErr)

		want := Err{Message: "failed to update personal deduction"}

		if !reflect.DeepEqual(gotErr, want) {
			t.Errorf("expected %v but got %v", want, gotErr)
		}

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %v but got %v", http.StatusInternalServerError, rec.Code)
		}
	})
}

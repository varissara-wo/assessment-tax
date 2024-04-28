package allowance

import (
	"errors"
	"testing"
)

func TestValidatePersoanlAmount(t *testing.T) {

	t.Run("should return nil if amount is valid", func(t *testing.T) {
		a := Amount{Amount: 20000.0}

		got := a.ValidatePersonal()

		if got != nil {
			t.Errorf("expected nil but got %v", got)
		}
	})

	t.Run("should return error if amount is less than 10000", func(t *testing.T) {
		a := Amount{Amount: 9999.0}

		want := errors.New(ErrInvalidPersonalGreaterAmount)
		got := a.ValidatePersonal()

		if got.Error() != want.Error() {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("should return error if amount is greater than 100000", func(t *testing.T) {
		a := Amount{Amount: 100001.0}

		want := errors.New(ErrInvalidPersonalLessAmount)
		got := a.ValidatePersonal()

		if got.Error() != want.Error() {
			t.Errorf("expected %v but got %v", want, got)
		}
	})
}

func TestValidateKreceiptAmount(t *testing.T) {
	t.Run("should return error if amount is less than 0", func(t *testing.T) {
		a := Amount{Amount: -1.0}

		want := errors.New(ErrInvalidKReceiptGreaterAmount)
		got := a.ValidateKReceipt()

		if got.Error() != want.Error() {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("should return error if amount is grether than 100000", func(t *testing.T) {
		a := Amount{Amount: 200000.0}

		want := errors.New(ErrInvalidKReceiptLessAmount)
		got := a.ValidateKReceipt()

		if got.Error() != want.Error() {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("should return nil if amount is valid", func(t *testing.T) {
		a := Amount{Amount: 20000.0}

		got := a.ValidateKReceipt()

		if got != nil {
			t.Errorf("expected nil but got %v", got)
		}
	})

}

func TestValidateAllowanceType(t *testing.T) {
	t.Run("should return error if allowance type is invalid", func(t *testing.T) {
		a := Allowance{AllowanceType: "invalid_type"}

		want := errors.New(ErrInvalidAllowance)
		got := validateAllowanceType(a)

		if got.Error() != want.Error() {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("should return nil if allowance type is valid", func(t *testing.T) {
		a := Allowance{AllowanceType: Donation}

		got := validateAllowanceType(a)

		if got != nil {
			t.Errorf("expected nil but got %v", got)
		}
	})

}

func TestValidateAllowanceAmount(t *testing.T) {
	t.Run("should return error if allowance amount is less than 0", func(t *testing.T) {
		a := Allowance{AllowanceType: Donation, Amount: -1.0}

		want := errors.New(ErrInvalidAllowanceAmount)
		got := ValidateAllowance(a)

		if got.Error() != want.Error() {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("should return nil if allowance amount is valid", func(t *testing.T) {
		a := Allowance{AllowanceType: Donation, Amount: 20000.0}

		got := ValidateAllowance(a)

		if got != nil {
			t.Errorf("expected nil but got %v", got)
		}
	})
}

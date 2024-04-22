package admin

import (
	"errors"
	"testing"
)

func TestValidatePersoanlAmount(t *testing.T) {

	t.Run("should return nil if amount is valid", func(t *testing.T) {
		a := Amount{Amount: 20000.0}

		got := a.ValidatePersonalDeduction()

		if got != nil {
			t.Errorf("expected nil but got %v", got)
		}
	})

	t.Run("should return error if amount is less than 10000", func(t *testing.T) {
		a := Amount{Amount: 9999.0}

		want := errors.New(ErrInvalidPersonalGreaterAmount)
		got := a.ValidatePersonalDeduction()

		if got.Error() != want.Error() {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("should return error if amount is greater than 100000", func(t *testing.T) {
		a := Amount{Amount: 100001.0}

		want := errors.New(ErrInvalidPersonalLessAmount)
		got := a.ValidatePersonalDeduction()

		if got.Error() != want.Error() {
			t.Errorf("expected %v but got %v", want, got)
		}
	})
}

package tax

import (
	"encoding/csv"
	"reflect"
	"strings"
	"testing"
)

func TestReadCSV(t *testing.T) {
	t.Run("should return TaxDetails slice and no error if CSV is valid", func(t *testing.T) {
		csvData := `totalIncome,wht,donation
1000.0,200.0,300.0
4000.0,500.0,600.0
`
		reader := csv.NewReader(strings.NewReader(csvData))

		got, err := readCSV(reader)

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		want := []TaxDetails{
			{
				TotalIncome: 1000.0,
				WHT:         200.0,
				Allowances: []Allowance{
					{
						AllowanceType: Donation,
						Amount:        300.0,
					},
				},
			},
			{
				TotalIncome: 4000.0,
				WHT:         500.0,
				Allowances: []Allowance{
					{
						AllowanceType: Donation,
						Amount:        600.0,
					},
				},
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("should return error if CSV header is invalid", func(t *testing.T) {
		csvData := `invalidHeader1,invalidHeader2,invalidHeader3
1000.0,200.0,300.0
`
		reader := csv.NewReader(strings.NewReader(csvData))

		_, got := readCSV(reader)

		want := ErrInvalidHeaderCSVData
		if got == nil || got.Error() != want {
			t.Errorf("expected error message %v but got %v", want, got)
		}
	})

	t.Run("should return error if CSV data value is empty", func(t *testing.T) {
		csvData := `totalIncome,wht,donation
1000.0,200.0,
`
		reader := csv.NewReader(strings.NewReader(csvData))

		_, got := readCSV(reader)

		want := ErrorInvalidEmptyCSVData
		if got == nil || got.Error() != want {
			t.Errorf("expected error message %v but got %v", want, got)
		}
	})
}

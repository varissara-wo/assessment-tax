package tax

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/varissara-wo/assessment-tax/allowance"
)

const (
	ErrInvalidHeaderCSVData  = "invalid CSV header, expected totalIncome, wht, donation"
	ErrorInvalidEmptyCSVData = "invalid CSV data value cannot be empty"
)

func readCSV(reader *csv.Reader) ([]TaxDetails, error) {
	row, err := reader.Read()
	if err != nil {
		return nil, err
	}

	if row[0] != "totalIncome" || row[1] != "wht" || row[2] != "donation" {
		return nil, errors.New(ErrInvalidHeaderCSVData)
	}

	tds := []TaxDetails{}
	for {
		row, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		td := TaxDetails{}

		if row[0] == "" || row[1] == "" || row[2] == "" {
			return nil, errors.New(ErrorInvalidEmptyCSVData)
		}

		for i, r := range row {
			v, err := strconv.ParseFloat(strings.Replace(r, ",", "", -1), 64)
			if err != nil {
				return nil, err
			}

			switch i {
			case 0:
				td.TotalIncome = v
			case 1:
				td.WHT = v
			case 2:
				td.Allowances = []allowance.Allowance{{
					AllowanceType: allowance.Donation,
					Amount:        v,
				}}
			}
		}

		tds = append(tds, td)

	}

	return tds, nil
}

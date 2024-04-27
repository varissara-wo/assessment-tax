package tax

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type Err struct {
	Message string `json:"message"`
}

type Storer interface {
	TaxCalculation(TaxDetails) (TaxCalculationResponse, error)
}

type Handler struct {
	store Storer
}

type Taxes struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
}

type TaxesResponse struct {
	Taxes []Taxes `json:"taxes"`
}

func New(store Storer) *Handler {
	return &Handler{store: store}
}

func (h *Handler) TaxHandler(c echo.Context) error {
	td := TaxDetails{}

	if err := c.Bind(&td); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if err := td.ValidateTaxDetails(); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	t, err := h.store.TaxCalculation(td)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, t)
}

func (h *Handler) TaxCSVHandler(c echo.Context) error {
	file, err := c.FormFile("file")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	defer src.Close()

	reader := csv.NewReader(src)
	// Read first line
	row, err := reader.Read()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	// Check header
	if row[0] != "totalIncome" || row[1] != "wht" || row[2] != "donation" {
		return c.JSON(http.StatusBadRequest, Err{Message: "Invalid CSV header"})
	}

	taxDetails := []TaxDetails{}
	for {
		row, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}

		td := TaxDetails{}

		if row[0] == "" || row[1] == "" || row[2] == "" {
			return c.JSON(http.StatusBadRequest, Err{Message: "Invalid CSV data value cannot be empty"})
		}

		totalIncome, err := strconv.ParseFloat(strings.Replace(row[0], ",", "", -1), 64)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}
		td.TotalIncome = totalIncome

		wht, err := strconv.ParseFloat(strings.Replace(row[1], ",", "", -1), 64)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}
		td.WHT = wht

		a, err := strconv.ParseFloat(strings.Replace(row[2], ",", "", -1), 64)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}
		td.Allowances = []Allowance{{
			AllowanceType: Donation,
			Amount:        a,
		}}

		taxDetails = append(taxDetails, td)

	}

	taxes := []Taxes{}
	// check
	for _, td := range taxDetails {

		if err := td.ValidateTaxDetails(); err != nil {
			return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
		}

		t, err := h.store.TaxCalculation(td)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}

		tr := Taxes{
			TotalIncome: td.TotalIncome,
			Tax:         t.Tax,
		}

		taxes = append(taxes, tr)

		log.Printf("Tax Detail: %+v\n", td)
	}

	taxesResponse := TaxesResponse{
		Taxes: taxes,
	}

	return c.JSON(http.StatusOK, taxesResponse)
}

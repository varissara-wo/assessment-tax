package tax

import (
	"encoding/csv"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Err struct {
	Message string `json:"message"`
}

type Storer interface {
	TaxCalculation(TaxDetails) (TaxResponse, error)
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
	taxDetails, err := readCSV(reader)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	taxes := []Taxes{}

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

	}

	taxesResponse := TaxesResponse{
		Taxes: taxes,
	}

	return c.JSON(http.StatusOK, taxesResponse)
}

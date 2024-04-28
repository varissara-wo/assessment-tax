package allowance

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Storer interface {
	SetPersonal(amount float64) (PersonalDeduction, error)
	SetKReceipt(amount float64) (KReceiptDeduction, error)
}

type Handler struct {
	store Storer
}

func New(store Storer) *Handler {
	return &Handler{store: store}
}

type Err struct {
	Message string `json:"message"`
}

func (h *Handler) SetPersonalHandler(c echo.Context) error {
	a := Amount{}

	if err := c.Bind(&a); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if err := a.ValidatePersonal(); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	p, err := h.store.SetPersonal(a.Amount)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, p)
}

func (h *Handler) SetKReceiptHandler(c echo.Context) error {
	a := Amount{}

	if err := c.Bind(&a); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if err := a.ValidateKReceipt(); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	k, err := h.store.SetKReceipt(a.Amount)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, k)
}

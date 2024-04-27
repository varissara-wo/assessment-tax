package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Storer interface {
	SetPersonalAllowance(amount float64) (PersonalAllowance, error)
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

type Amount struct {
	Amount float64 `json:"amount"`
}

type PersonalAllowance struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

func (h *Handler) SetPersonalAllowanceHandler(c echo.Context) error {
	a := Amount{}

	if err := c.Bind(&a); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if err := a.ValidatePersonalAllowance(); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	d, err := h.store.SetPersonalAllowance(a.Amount)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, d)
}

package allowance

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Storer interface {
	SetKReceipt(amount float64) (KReceipt, error)
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

package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/varissara-wo/assessment-tax/admin"
	"github.com/varissara-wo/assessment-tax/postgres"
	"github.com/varissara-wo/assessment-tax/tax"
)

func main() {
	p, err := postgres.New()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	th := tax.New(p)
	ah := admin.New(p)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	e.POST("/tax/calculations", th.TaxHandler)
	e.PUT("/admin/personal-deduction", ah.PersonalDeductionHandler)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

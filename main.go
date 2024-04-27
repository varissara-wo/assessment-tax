package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
	e.POST("/tax/calculations", th.TaxHandler)
	e.POST("tax/calculations/upload-csv", th.TaxCSVHandler)

	ah := admin.New(p)
	a := e.Group("/admin")
	a.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == os.Getenv("ADMIN_USERNAME") && password == os.Getenv("ADMIN_PASSWORD") {
			return true, nil
		}
		return false, nil

	}))

	a.POST("/personal-deduction", ah.SetPersonalAllowanceHandler)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

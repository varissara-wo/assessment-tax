package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/varissara-wo/assessment-tax/allowance"
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
	e.POST("/tax/calculations/upload-csv", th.TaxCSVHandler)

	aw := allowance.New(p)
	a := e.Group("/admin")
	a.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == os.Getenv("ADMIN_USERNAME") && password == os.Getenv("ADMIN_PASSWORD") {
			return true, nil
		}
		return false, nil

	}))

	a.POST("/deductions/personal", aw.SetPersonalHandler)
	a.POST("/deductions/k-receipt", aw.SetKReceiptHandler)

	go func() {
		if err := e.Start(":" + os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Info("start error, shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	log.Println("shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

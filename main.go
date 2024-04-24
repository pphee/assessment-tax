package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pphee/assessment-tax/module/tax"
	"github.com/pphee/assessment-tax/store"
	"net/http"
	"os"
)

func main() {
	e := echo.New()

	// Database connection setup
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		e.Logger.Fatal("DATABASE_URL not set in environment variables")
	}
	dbStore := store.NewPostgresStore(dsn)

	sqlDB, err := dbStore.DB.DB()
	if err != nil {
		e.Logger.Fatal("Failed to get raw database object: ", err)
	}
	defer sqlDB.Close()

	// Basic Auth for admin routes
	adminUser := os.Getenv("ADMIN_USERNAME")
	adminPass := os.Getenv("ADMIN_PASSWORD")
	if adminUser == "" || adminPass == "" {
		e.Logger.Fatal("Admin username or password not set in environment variables")
	}

	taxRepo := tax.NewTaxRepository(dbStore.DB)
	taxService := tax.NewTaxService(taxRepo)
	taxHandler := tax.NewTaxHandler(taxService)

	taxGroup := e.Group("/tax")
	taxGroup.POST("/calculations", taxHandler.PostTaxCalculation)
	taxGroup.POST("/calculations/upload-csv", taxHandler.TaxCalculationsCSVHandler)

	admin := e.Group("/admin")
	admin.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		return username == adminUser && password == adminPass, nil
	}))
	{
		admin.POST("/deductions/personal", taxHandler.SetPersonalDeduction)
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}

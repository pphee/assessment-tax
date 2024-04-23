package main

import (
	"github.com/labstack/echo/v4"
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

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

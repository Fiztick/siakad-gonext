package main

import (
	"flag"
	"net/http"
	"siakad-backend/pkg/database"

	"github.com/labstack/echo/v5"
)

func main() {
	shouldSeed := flag.Bool("seed", false, "Set to true to seed the database")
	flag.Parse()

	db := database.InitDB()
	e := echo.New()
	// h := &handler.Handler{DB: db}
	if *shouldSeed {
		database.Seed(db)
	}

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "SIAKAD API is Online",
			"db":     "Connected",
		})
	})

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}

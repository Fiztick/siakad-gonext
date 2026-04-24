package main

import (
	"flag"
	"net/http"
	"os"
	"siakad-backend/internal/handler"
	"siakad-backend/pkg/database"

	echoJWT "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func main() {
	shouldSeed := flag.Bool("seed", false, "Set to true to seed the database")
	flag.Parse()

	db := database.InitDB()
	e := echo.New()
	h := &handler.Handler{DB: db}
	if *shouldSeed {
		database.Seed(db)
		if flag.Lookup("seed").Value.String() == "true" {
			return
		}
	}

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "SIAKAD API is Online",
			"db":     "Connected",
		})
	})
	e.POST("/login", h.Login)

	// Protected Routes
	protected := e.Group("/api")
	protected.Use(echoJWT.WithConfig(echoJWT.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	protected.POST("/students", h.CreateStudent)
	protected.GET("/student/:studentId", h.GetStudent)
	protected.PUT("/student/:studentId", h.UpdateStudent)
	protected.DELETE("/student/:studentId", h.DeleteStudent)

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}

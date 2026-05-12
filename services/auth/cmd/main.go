package main

import (
	"net/http"

	"github.com/Anicet78/SolanumStreaming/auth/internal/handler"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger())

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	authHandler := handler.NewAuthHandler()
	authHandler.RegisterRoutes(e)

	if err := e.Start(":8081"); err != nil {
		e.Logger.Error("Failed to start server", "error", err)
	}
}

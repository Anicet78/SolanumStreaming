package main

import (
	"net/http"
	"shared/database"

	"github.com/Anicet78/SolanumStreaming/auth/internal/handler"
	"github.com/Anicet78/SolanumStreaming/auth/internal/service"
	"github.com/Anicet78/SolanumStreaming/auth/internal/store"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	pool, err := database.Connect()
	if err != nil {
		return
	}
	defer pool.Close()

	store := store.New(pool)

	authService := service.NewAuthService(store)

	authHandler := handler.NewAuthHandler(authService)

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Validator = &CustomValidator{validator: validator.New()}

	authHandler.RegisterRoutes(e)

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	if err := e.Start(":8081"); err != nil {
		e.Logger.Error("Failed to start server", "error", err)
	}
}

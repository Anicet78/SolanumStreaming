package main

import (
	"net/http"

	"github.com/Anicet78/SolanumStreaming/stream/internal/handler"
	"github.com/Anicet78/SolanumStreaming/stream/internal/service"
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
	streamService := service.NewStreamService()

	streamHandler := handler.NewStreamHandler(streamService)

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Validator = &CustomValidator{validator: validator.New()}

	streamHandler.RegisterRoutes(e)

	e.RouteNotFound("/*", func(c *echo.Context) error {
		return echo.NewHTTPError(http.StatusNotFound, "Route Not Found")
	})

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	if err := e.Start(":8082"); err != nil {
		e.Logger.Error("Failed to start server", "error", err)
	}
}

package main

import (
	"net/http"
	"shared/database"

	"github.com/Anicet78/SolanumStreaming/movies/internal/handler"
	"github.com/Anicet78/SolanumStreaming/movies/internal/service"
	"github.com/Anicet78/SolanumStreaming/movies/internal/store"
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

	movieService := service.NewMovieService(store)

	movieHandler := handler.NewMovieHandler(movieService)

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5500"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))
	e.Use(middleware.RequestLogger())
	e.Validator = &CustomValidator{validator: validator.New()}

	movieHandler.RegisterRoutes(e)

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

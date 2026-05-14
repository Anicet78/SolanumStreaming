package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Anicet78/SolanumStreaming/auth/internal/handler"
	"github.com/Anicet78/SolanumStreaming/auth/internal/service"
	"github.com/Anicet78/SolanumStreaming/auth/internal/store"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func ConnectToDB() (*pgxpool.Pool, error) {
	var ctx = context.Background()

	// Initialize the connection pool
	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		return nil, err
	}

	// Verify the connection
	if err := pool.Ping(ctx); err != nil {
		slog.Error("Unable to ping database", "error", err)
		return nil, err
	}

	fmt.Println("Connected to PostgreSQL database!")
	return pool, nil
}

func main() {
	pool, err := ConnectToDB()
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

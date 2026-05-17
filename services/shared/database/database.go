package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() (*pgxpool.Pool, error) {
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

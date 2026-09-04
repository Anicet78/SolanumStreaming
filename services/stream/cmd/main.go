package main

import (
	"log/slog"
	"net/http"

	"github.com/Anicet78/SolanumStreaming/stream/internal/handler"
	"github.com/Anicet78/SolanumStreaming/stream/internal/service"
	"github.com/Anicet78/SolanumStreaming/stream/internal/store"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
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
	ramProvider := store.NewProvider()

	capacityFunc := func() (int64, bool) {
		return 512 * 1024 * 1024, true
	}
	var tc storage.TorrentCapacity = &capacityFunc

	ramStorage := storage.NewResourcePiecesOpts(ramProvider, storage.ResourcePiecesOpts{
		Capacity: tc,
	})

	cfg := torrent.NewDefaultClientConfig()
	cfg.DefaultStorage = ramStorage

	client, err := torrent.NewClient(cfg)
	if err != nil {
		slog.Error("Unable create a torrent client", "error", err)
	}
	defer client.Close()

	streamService := service.NewStreamService(client)

	streamHandler := handler.NewStreamHandler(streamService)

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5500"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))
	e.Use(middleware.RequestLogger())
	e.Validator = &CustomValidator{validator: validator.New()}

	streamHandler.RegisterRoutes(e)

	e.RouteNotFound("/*", func(c *echo.Context) error {
		return echo.NewHTTPError(http.StatusNotFound, "Route Not Found")
	})

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	if err := e.Start(":8083"); err != nil {
		e.Logger.Error("Failed to start server", "error", err)
	}
}

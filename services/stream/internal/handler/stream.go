package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"shared/bind"
	"shared/response"
	"time"

	"github.com/Anicet78/SolanumStreaming/stream/internal/domain"
	"github.com/Anicet78/SolanumStreaming/stream/internal/service"
	"github.com/labstack/echo/v5"
)

type StreamHandler struct {
	service *service.StreamService
}

func NewStreamHandler(service *service.StreamService) *StreamHandler {
	return &StreamHandler{service: service}
}

func (h *StreamHandler) RegisterRoutes(e *echo.Echo) {
	private := e.Group("")
	// private.Use(auth.JWTMiddleware())
	private.GET("/", h.stream)
}

func (h *StreamHandler) stream(c *echo.Context) error {
	params, err := bind.Query[domain.TorrentLinkQuery](c)
	if err != nil {
		return err
	}

	reader, err := h.service.Stream(c.Request().Context(), params.TorrentLink)
	if err != nil {
		if errors.Is(err, domain.ErrTorrentLoadingTimeout) {
			return response.GatewayTimeout(c, err.Error())
		}
		slog.Error("Torrent stream failed", "error", err)
		return response.InternalServerError(c, "internal server error")
	}

	defer reader.Close()
	http.ServeContent(c.Response(), c.Request(), "test", time.Time{}, reader)

	return response.NoContent(c)
}

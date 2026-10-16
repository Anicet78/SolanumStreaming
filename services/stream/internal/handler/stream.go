package handler

import (
	"errors"
	"log/slog"
	"shared/bind"
	"shared/response"

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
	private.GET("/playlist.m3u8", h.serveRaw)
	// private.GET("/:segment", h.serveFile)
}

func (h *StreamHandler) serveRaw(c *echo.Context) error {
	params, err := bind.Query[domain.TorrentLinkQuery](c)
	if err != nil {
		return err
	}

	err = h.service.ServeRaw(c.Response(), c.Request(), params.TorrentLink)
	if err != nil {
		if errors.Is(err, domain.ErrTorrentLoadingTimeout) {
			return response.GatewayTimeout(c, err.Error())
		} else if errors.Is(err, domain.ErrTorrentNoFile) {
			return response.ServiceUnavailable(c, err.Error())
		}
		slog.Error("Torrent stream failed", "error", err)
		return response.InternalServerError(c, "internal server error")
	}

	return response.NoContent(c)
}

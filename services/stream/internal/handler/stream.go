package handler

import (
	"errors"
	"log/slog"
	"shared/auth"
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
	return &StreamHandler{
		service: service,
	}
}

func (h *StreamHandler) RegisterRoutes(e *echo.Echo) {
	private := e.Group("")
	// private.Use(auth.JWTMiddleware())
	private.GET("/playlist.m3u8", h.createPlaylist)
	// private.GET("/:segment", h.serveFile)

	internal := e.Group("")
	internal.Use(auth.InternalMiddleware(h.service.InternalToken))
	internal.GET("/raw", h.raw)
}

func (h *StreamHandler) createPlaylist(c *echo.Context) error {
	params, err := bind.Query[domain.TorrentLinkQuery](c)
	if err != nil {
		return err
	}

	res, err := h.service.CreatePlaylist(c.Request().Context(), params.TorrentLink)
	if err != nil {
		if errors.Is(err, domain.ErrTorrentLoadingTimeout) {
			return response.GatewayTimeout(c, err.Error())
		} else if errors.Is(err, domain.ErrTorrentNoFile) {
			return response.ServiceUnavailable(c, err.Error())
		} else if errors.Is(err, domain.ErrNoVideo) {
			return response.ServiceUnavailable(c, err.Error())
		} else if errors.Is(err, domain.ErrUnsupportedContainer) {
			return response.UnsupportedMediaType(c, err.Error())
		}
		slog.Error("Torrent stream failed", "error", err)
		return response.InternalServerError(c, "internal server error")
	}

	return response.OK(c, res)
}

func (h *StreamHandler) raw(c *echo.Context) error {
	return h.service.ServeRaw(c.Response(), c.Request(), c.QueryParam("link"))
}

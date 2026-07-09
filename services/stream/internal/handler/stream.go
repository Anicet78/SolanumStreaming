package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"shared/auth"
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
	private.Use(auth.JWTMiddleware())
	private.GET("/:torrent_link", h.stream)
}

func (h *StreamHandler) stream(c *echo.Context) error {
	params, err := bind.Params[domain.TorrentLinkParam](c)
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

	http.ServeContent(c.Response(), c.Request(), "test", time.Time{}, reader)
	defer reader.Close()

	return response.NoContent(c)
}

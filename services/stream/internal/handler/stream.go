package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"shared/auth"
	"shared/bind"
	"shared/response"
	"strings"

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
	private.GET("/hls/:id/:file", h.hlsFile)

	internal := e.Group("")
	internal.Use(auth.InternalMiddleware(h.service.InternalToken))
	internal.GET("/raw", h.raw)
}

func (h *StreamHandler) createPlaylist(c *echo.Context) error {
	params, err := bind.Query[domain.TorrentLinkQuery](c)
	if err != nil {
		return err
	}

	id, err := h.service.CreatePlaylist(c.Request().Context(), params.TorrentLink)
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

	return response.Found(c, "/hls/"+id+"/index.m3u8")
}

func (h *StreamHandler) hlsFile(c *echo.Context) error {
	id := c.Param("id")
	file := c.Param("file")

	pl, ok := h.service.GetPlaylist(id)
	if !ok {
		return response.NotFound(c, "Playlist not found")
	}

	switch {
	case file == "index.m3u8":
		c.Response().Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		c.Response().Header().Set("Cache-Control", "no-cache")
		return c.String(http.StatusOK, pl.Text)
	case file == "init.mp4":
		data, err := h.service.GetInit(c.Request().Context(), pl)
		if err != nil {
			return h.mediaError(c, "init", err)
		}
		c.Response().Header().Set("Cache-Control", "public, max-age=3600")
		return c.Blob(http.StatusOK, "video/mp4", data)
	case strings.HasPrefix(file, "seg_") && strings.HasSuffix(file, ".m4s"):
		var idx int
		if _, err := fmt.Sscanf(file, "seg_%d.m4s", &idx); err != nil {
			return response.NotFound(c, "Not found")
		}
		data, err := h.service.GetSegment(c.Request().Context(), pl, idx)
		if err != nil {
			return h.mediaError(c, fmt.Sprintf("segment %d", idx), err)
		}
		c.Response().Header().Set("Cache-Control", "public, max-age=3600")
		return c.Blob(http.StatusOK, "video/mp4", data)
	default:
		return response.NotFound(c, "Not found")
	}
}

func (h *StreamHandler) mediaError(c *echo.Context, what string, err error) error {
	switch {
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return nil
	case errors.Is(err, domain.ErrSegmentOutOfRange):
		return response.NotFound(c, "Segment not found")
	}
	slog.Error("media request failed", "what", what, "error", err)
	return response.InternalServerError(c, "Internal server error")
}

func (h *StreamHandler) raw(c *echo.Context) error {
	return h.service.ServeRaw(c.Response(), c.Request(), c.QueryParam("link"))
}

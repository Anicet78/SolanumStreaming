package handler

import (
	"log/slog"
	"shared/auth"
	"shared/response"

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
	private.GET("/:stream_id", h.stream)
}

func (h *StreamHandler) stream(c *echo.Context) error {
	res, err := h.service.Stream(c.Request().Context())
	if err != nil {
		slog.Error("Torrent stream failed", "error", err)
		return response.InternalServerError(c, "internal server error")
	}

	return response.OK(c, res)
}

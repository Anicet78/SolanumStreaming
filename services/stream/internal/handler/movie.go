package handler

import (
	"log/slog"
	"shared/auth"
	"shared/bind"
	"shared/response"

	"github.com/Anicet78/SolanumStreaming/stream/internal/domain"
	"github.com/Anicet78/SolanumStreaming/stream/internal/service"
	"github.com/labstack/echo/v5"
)

type MovieHandler struct {
	service *service.MovieService
}

func NewMovieHandler(service *service.MovieService) *MovieHandler {
	return &MovieHandler{service: service}
}

func (h *MovieHandler) RegisterRoutes(e *echo.Echo) {
	private := e.Group("")
	private.Use(auth.JWTMiddleware())
	private.GET("/:movie_id", h.stream)
}

func (h *MovieHandler) stream(c *echo.Context) error {
	params, err := bind.Query[domain.](c)
	if err != nil {
		return err
	}

	res, err := h.service.Stream(c.Request().Context(), params)
	if err != nil {
		slog.Error("Torrent stream failed", "error", err)
		return response.InternalServerError(c, "internal server error")
	}

	return response.OK(c, res)
}

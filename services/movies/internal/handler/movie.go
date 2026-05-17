package handler

import (
	"log/slog"
	"shared/auth"
	"shared/bind"
	"shared/response"

	"github.com/Anicet78/SolanumStreaming/movies/internal/domain"
	"github.com/Anicet78/SolanumStreaming/movies/internal/service"
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
	private.POST("/search", h.search)
}

func (h *MovieHandler) search(c *echo.Context) error {
	_, err := bind.Body[domain.SearchRequestBody](c)
	if err != nil {
		return err
	}

	res, err := h.service.Search(c.Request().Context())
	if err != nil {
		slog.Error("Search creation failed", "error", err)
		return response.InternalServerError(c, "internal server error")
	}

	return response.Created(c, res)
}

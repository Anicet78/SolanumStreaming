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
	private.GET("/search", h.search)
	private.POST("/add", h.add)
}

func (h *MovieHandler) search(c *echo.Context) error {
	params, err := bind.Query[domain.SearchRequestQuery](c)
	if err != nil {
		return err
	}

	res, err := h.service.Search(c.Request().Context(), params)
	if err != nil {
		slog.Error("Search creation failed", "error", err)
		return response.InternalServerError(c, "internal server error")
	}

	return response.OK(c, res)
}

func (h *MovieHandler) add(c *echo.Context) error {
	body, err := bind.Body[domain.AddRequestBody](c)
	if err != nil {
		return err
	}

	uuid, err := auth.StringToUUID(c.Get("uuid").(string))
	if err != nil {
		return response.InternalServerError(c, "internal server error")
	}

	err = h.service.Add(c.Request().Context(), uuid, body.MovieID)
	if err != nil {
		slog.Error("Add movie to collection failed", "error", err)
		return response.InternalServerError(c, "internal server error")
	}

	return response.NoContent(c)
}

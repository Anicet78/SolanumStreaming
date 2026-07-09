package handler

import (
	"errors"
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
	private.GET("/collection", h.getCollection)
	private.GET("/collection/:movie_id", h.getInCollection)
	private.POST("/collection", h.addToCollection)
	private.DELETE("/collection/:movie_id", h.removeFromCollection)
}

func (h *MovieHandler) search(c *echo.Context) error {
	params, err := bind.Query[domain.SearchRequestQuery](c)
	if err != nil {
		return err
	}

	res, err := h.service.Search(c.Request().Context(), params)
	if err != nil {
		slog.Error("Search failed", "error", err)
		return response.InternalServerError(c, "internal server error")
	}

	return response.OK(c, res)
}

func (h *MovieHandler) getCollection(c *echo.Context) error {
	uuid, err := auth.StringToUUID(c.Get("uuid").(string))
	if err != nil {
		return response.InternalServerError(c, "internal server error")
	}

	res, err := h.service.GetCollection(c.Request().Context(), uuid)
	if err != nil {
		slog.Error("Get collection failed", "error", err)
		return response.InternalServerError(c, "internal server error")
	}

	return response.OK(c, res)
}

func (h *MovieHandler) getInCollection(c *echo.Context) error {
	params, err := bind.Params[domain.MovieIDParam](c)
	if err != nil {
		return err
	}

	uuid, err := auth.StringToUUID(c.Get("uuid").(string))
	if err != nil {
		return response.InternalServerError(c, "internal server error")
	}

	res, err := h.service.GetInCollection(c.Request().Context(), uuid, params.MovieID)
	if err != nil {
		if errors.Is(err, domain.ErrMovieNotInCollection) {
			return response.NotFound(c, err.Error())
		}
		slog.Error("Get movie in collection failed", "error", err)
		return response.InternalServerError(c, "internal server error")
	}

	return response.OK(c, res)
}

func (h *MovieHandler) addToCollection(c *echo.Context) error {
	body, err := bind.Body[domain.CollectionRequestBody](c)
	if err != nil {
		return err
	}

	uuid, err := auth.StringToUUID(c.Get("uuid").(string))
	if err != nil {
		return response.InternalServerError(c, "internal server error")
	}

	res, err := h.service.AddToCollection(c.Request().Context(), uuid, body.MovieID)
	if err != nil {
		if errors.Is(err, domain.ErrMovieAlreadyInCollection) {
			return response.Conflict(c, err.Error())
		}
		slog.Error("Add movie to collection failed", "error", err)
		return response.InternalServerError(c, "internal server error")
	}

	return response.OK(c, res)
}

func (h *MovieHandler) removeFromCollection(c *echo.Context) error {
	params, err := bind.Params[domain.MovieIDParam](c)
	if err != nil {
		return err
	}

	uuid, err := auth.StringToUUID(c.Get("uuid").(string))
	if err != nil {
		return response.InternalServerError(c, "internal server error")
	}

	err = h.service.RemoveFromCollection(c.Request().Context(), uuid, params.MovieID)
	if err != nil {
		if errors.Is(err, domain.ErrMovieNotInCollection) {
			return response.NotFound(c, err.Error())
		}
		slog.Error("Remove movie from collection failed", "error", err)
		return response.InternalServerError(c, "internal server error")
	}

	return response.NoContent(c)
}

package handler

import (
	"errors"
	"log/slog"
	"shared/auth"
	"shared/bind"
	"shared/response"

	"github.com/Anicet78/SolanumStreaming/auth/internal/domain"
	"github.com/Anicet78/SolanumStreaming/auth/internal/service"
	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) RegisterRoutes(e *echo.Echo) {
	public := e.Group("")
	public.POST("/register", h.register)
	public.POST("/login", h.login)

	private := e.Group("")
	private.Use(auth.JWTMiddleware())
	private.DELETE("/profile", h.delete)
	private.PATCH("/profile", h.patchProfile)
}

func (h *AuthHandler) register(c *echo.Context) error {
	req, err := bind.Body[domain.CreateUserRequest](c)
	if err != nil {
		return err
	}

	res, err := h.service.Register(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUsernameAlreadyExists) {
			return response.Conflict(c, err.Error())
		}
		slog.Error("Account creation failed", "error", err)
		return response.InternalServerError(c, "internal server error")
	}

	return response.Created(c, res)
}

func (h *AuthHandler) login(c *echo.Context) error {
	req, err := bind.Body[domain.LoginUserRequest](c)
	if err != nil {
		return err
	}

	res, err := h.service.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUsernameDoesNotExists) {
			return response.NotFound(c, err.Error())
		} else if errors.Is(err, domain.ErrPasswordDoesNotMatch) {
			return response.BadRequest(c, err.Error())
		}
		slog.Error("Log into account failed", "error", err)
		return response.InternalServerError(c, "internal server error")
	}

	return response.OK(c, res)
}

func (h *AuthHandler) delete(c *echo.Context) error {
	uuid, err := auth.StringToUUID(c.Get("uuid").(string))
	if err != nil {
		return response.InternalServerError(c, "internal server error")
	}

	err = h.service.Delete(c.Request().Context(), uuid)
	if err != nil {
		slog.Error("Account deletion failed", "error", err)
		return response.InternalServerError(c, "internal server error")
	}
	return response.NoContent(c)
}

func (h *AuthHandler) patchProfile(c *echo.Context) error {
	req, err := bind.Body[domain.PatchProfileRequest](c)
	if err != nil {
		return err
	}

	uuid, err := auth.StringToUUID(c.Get("uuid").(string))
	if err != nil {
		return response.InternalServerError(c, "internal server error")
	}

	err = h.service.PatchProfile(c.Request().Context(), uuid, req.NewUsername)
	if err != nil {
		if errors.Is(err, domain.ErrUsernameAlreadyExists) {
			return response.Conflict(c, err.Error())
		}
		slog.Error("Cannot update profile", "error", err)
		return response.InternalServerError(c, "internal server error")
	}

	return response.NoContent(c)
}

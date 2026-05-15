package handler

import (
	"errors"
	"log/slog"
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
	e.POST("/register", h.register)
	e.POST("/login", h.login)
}

func (h *AuthHandler) register(c *echo.Context) error {
	var req domain.CreateUserRequest

	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid body")
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, err.Error())
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
	var req domain.LoginUserRequest

	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid body")
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, err.Error())
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

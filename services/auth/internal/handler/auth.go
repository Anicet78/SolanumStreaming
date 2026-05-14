package handler

import (
	"shared/response"

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
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=4"`
}

func (h *AuthHandler) register(c *echo.Context) error {
	var req RegisterRequest

	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid body")
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, req)
}

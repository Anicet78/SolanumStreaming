package handler

import (
	"shared/response"

	"github.com/labstack/echo/v5"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/register", h.register)
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) register(c *echo.Context) error {
	var req RegisterRequest

	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid body")
	}

	return response.OK(c, req)
}

package handler

import (
	"net/http"
	"github.com/labstack/echo/v5"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/register", h.register)
	e.POST("/login", h.login)
}

func (h *AuthHandler) register(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "register",
	})
}

func (h *AuthHandler) login(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "login",
	})
}

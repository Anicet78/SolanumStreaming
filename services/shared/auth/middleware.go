package auth

import (
	"shared/response"
	"strings"

	"github.com/labstack/echo/v5"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				return response.Unauthorized(c, "missing token")
			}

			token := strings.TrimPrefix(auth, "Bearer ")

			claims, err := VerifyToken(token)
			if err != nil {
				return response.Unauthorized(c, "invalid token")
			}

			c.Set("uuid", claims.UUID)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}

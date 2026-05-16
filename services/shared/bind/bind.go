package bind

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func Body[T any](c *echo.Context) (T, error) {
	var req T

	if err := c.Bind(&req); err != nil {
		return req, echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	if err := c.Validate(&req); err != nil {
		return req, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return req, nil
}

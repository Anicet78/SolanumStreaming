package bind

import (
	"shared/response"

	"github.com/labstack/echo/v5"
)

func Query[T any](c *echo.Context) (T, error) {
	var req T

	if err := c.Bind(&req); err != nil {
		return req, response.BadRequest(c, "invalid query params")
	}

	if err := c.Validate(&req); err != nil {
		return req, response.BadRequest(c, err.Error())
	}

	return req, nil
}

func Body[T any](c *echo.Context) (T, error) {
	var req T

	if err := c.Bind(&req); err != nil {
		return req, response.BadRequest(c, "invalid body")
	}

	if err := c.Validate(&req); err != nil {
		return req, response.BadRequest(c, err.Error())
	}

	return req, nil
}

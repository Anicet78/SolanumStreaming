package bind

import (
	"shared/response"

	"github.com/labstack/echo/v5"
)

func Params[T any](c *echo.Context) (T, error) {
	var req T

	if err := c.Bind(&req); err != nil {
		response.BadRequest(c, "invalid path params")
		return req, err
	}

	if err := c.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return req, err
	}

	return req, nil
}

func Query[T any](c *echo.Context) (T, error) {
	var req T

	if err := c.Bind(&req); err != nil {
		response.BadRequest(c, "invalid query params")
		return req, err
	}

	if err := c.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return req, err
	}

	return req, nil
}

func Body[T any](c *echo.Context) (T, error) {
	var req T

	if err := c.Bind(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return req, err
	}

	if err := c.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return req, err
	}

	return req, nil
}

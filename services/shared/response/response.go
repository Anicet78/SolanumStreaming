package response

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Success struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func JSON(c *echo.Context, status int, data any) error {
	return c.JSON(status, Success{Data: data})
}

func Message(c *echo.Context, status int, message string) error {
	if message == "" {
		message = http.StatusText(status)
	}
	return c.JSON(status, Success{Message: message})
}

func MessageData(c *echo.Context, status int, message string, data any) error {
	if message == "" {
		message = http.StatusText(status)
	}
	return c.JSON(status, Success{Message: message, Data: data})
}

func Fail(c *echo.Context, status int, message string) error {
	if message == "" {
		message = http.StatusText(status)
	}
	return c.JSON(status, Error{Message: message, Code: status})
}

func Failf(c *echo.Context, status int, format string, args ...any) error {
	return Fail(c, status, fmt.Sprintf(format, args...))
}

// --- 2xx ---

func OK(c *echo.Context, data any) error {
	return JSON(c, http.StatusOK, data)
}

func OKMessage(c *echo.Context, message string) error {
	return Message(c, http.StatusOK, message)
}

func OKMessageData(c *echo.Context, message string, data any) error {
	return MessageData(c, http.StatusOK, message, data)
}

func Created(c *echo.Context, data any) error {
	return JSON(c, http.StatusCreated, data)
}

func Accepted(c *echo.Context, data any) error {
	return JSON(c, http.StatusAccepted, data)
}

func NoContent(c *echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// --- 3xx ---

func MovedPermanently(c *echo.Context, url string) error {
	return c.Redirect(http.StatusMovedPermanently, url)
}

func Found(c *echo.Context, url string) error {
	return c.Redirect(http.StatusFound, url)
}

func SeeOther(c *echo.Context, url string) error {
	return c.Redirect(http.StatusSeeOther, url)
}

func TemporaryRedirect(c *echo.Context, url string) error {
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func PermanentRedirect(c *echo.Context, url string) error {
	return c.Redirect(http.StatusPermanentRedirect, url)
}

// --- 4xx ---

func BadRequest(c *echo.Context, message string) error {
	return Fail(c, http.StatusBadRequest, message)
}

func Unauthorized(c *echo.Context, message string) error {
	return Fail(c, http.StatusUnauthorized, message)
}

func PaymentRequired(c *echo.Context, message string) error {
	return Fail(c, http.StatusPaymentRequired, message)
}

func Forbidden(c *echo.Context, message string) error {
	return Fail(c, http.StatusForbidden, message)
}

func NotFound(c *echo.Context, message string) error {
	return Fail(c, http.StatusNotFound, message)
}

func MethodNotAllowed(c *echo.Context, message string) error {
	return Fail(c, http.StatusMethodNotAllowed, message)
}

func NotAcceptable(c *echo.Context, message string) error {
	return Fail(c, http.StatusNotAcceptable, message)
}

func RequestTimeout(c *echo.Context, message string) error {
	return Fail(c, http.StatusRequestTimeout, message)
}

func Conflict(c *echo.Context, message string) error {
	return Fail(c, http.StatusConflict, message)
}

func Gone(c *echo.Context, message string) error {
	return Fail(c, http.StatusGone, message)
}

func UnsupportedMediaType(c *echo.Context, message string) error {
	return Fail(c, http.StatusUnsupportedMediaType, message)
}

func UnprocessableEntity(c *echo.Context, message string) error {
	return Fail(c, http.StatusUnprocessableEntity, message)
}

func TooManyRequests(c *echo.Context, message string) error {
	return Fail(c, http.StatusTooManyRequests, message)
}

// --- 5xx ---

func InternalServerError(c *echo.Context, message string) error {
	return Fail(c, http.StatusInternalServerError, message)
}

func NotImplemented(c *echo.Context, message string) error {
	return Fail(c, http.StatusNotImplemented, message)
}

func BadGateway(c *echo.Context, message string) error {
	return Fail(c, http.StatusBadGateway, message)
}

func ServiceUnavailable(c *echo.Context, message string) error {
	return Fail(c, http.StatusServiceUnavailable, message)
}

func GatewayTimeout(c *echo.Context, message string) error {
	return Fail(c, http.StatusGatewayTimeout, message)
}

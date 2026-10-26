package service

import (
	"net/http"
	"time"

	"github.com/Anicet78/SolanumStreaming/auth/internal/store"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v5"
)

func AddLangCookie(c *echo.Context, language store.Languages) {
	langCookie := new(http.Cookie)
	langCookie.Name = "lang"
	langCookie.Value = string(language)
	langCookie.Path = "/"
	c.SetCookie(langCookie)
}

func AddNewRefreshTokenCookie(c *echo.Context, s *store.Queries, userID pgtype.UUID) error {
	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = newRefreshToken
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.MaxAge = 7 * 24 * 60 * 60

	_, err = s.CreateRefreshToken(c.Request().Context(), store.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: newRefreshToken,
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().UTC().Add(time.Duration(cookie.MaxAge)),
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	c.SetCookie(cookie)
	return nil
}

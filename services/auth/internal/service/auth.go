package service

import (
	"context"
	"errors"
	"net/http"
	"shared/auth"
	"time"

	"github.com/Anicet78/SolanumStreaming/auth/internal/domain"
	"github.com/Anicet78/SolanumStreaming/auth/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v5"
)

type AuthService struct {
	store *store.Queries
}

func NewAuthService(store *store.Queries) *AuthService {
	return &AuthService{store: store}
}

func (s *AuthService) Register(c *echo.Context, username string, password string) (domain.CreateUserResponse, error) {
	exists, err := s.store.GetUserByUsername(c.Request().Context(), username)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return domain.CreateUserResponse{}, err
		}
	} else if exists.Uuid.Valid {
		return domain.CreateUserResponse{}, domain.ErrUsernameAlreadyExists
	}

	hash, err := hashPassword(password)
	if err != nil {
		return domain.CreateUserResponse{}, err
	}

	created, err := s.store.CreateUser(c.Request().Context(), store.CreateUserParams{
		Username: username,
		Password: hash,
	})
	if err != nil {
		return domain.CreateUserResponse{}, err
	}

	token, err := auth.GenerateToken(created.Uuid.String(), string(created.Role))
	if err != nil {
		return domain.CreateUserResponse{}, err
	}

	token, err = generateRefreshToken()
	if err != nil {
		return domain.CreateUserResponse{}, err
	}

	cookie := new(http.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.MaxAge = 7 * 24 * 60 * 60

	c.SetCookie(cookie)

	return domain.CreateUserResponse{
		UUID:     created.Uuid.String(),
		Username: created.Username,
		Role:     string(created.Role),
		JWT:      token,
	}, nil
}

func (s *AuthService) Login(c *echo.Context, username string, password string) (domain.LoginUserResponse, error) {
	found, err := s.store.GetUserByUsername(c.Request().Context(), username)
	if err != nil {
		return domain.LoginUserResponse{}, domain.ErrUsernameDoesNotExists
	}

	match, err := passwordMatch(password, found.Password)
	if err != nil {
		return domain.LoginUserResponse{}, err
	} else if !match {
		return domain.LoginUserResponse{}, domain.ErrPasswordDoesNotMatch
	}

	token, err := auth.GenerateToken(found.Uuid.String(), string(found.Role))
	if err != nil {
		return domain.LoginUserResponse{}, err
	}

	token, err = generateRefreshToken()
	if err != nil {
		return domain.LoginUserResponse{}, err
	}

	cookie := new(http.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.MaxAge = 7 * 24 * 60 * 60

	c.SetCookie(cookie)

	return domain.LoginUserResponse{
		UUID:     found.Uuid.String(),
		Username: found.Username,
		Role:     string(found.Role),
		JWT:      token,
	}, nil
}

func (s *AuthService) Refresh(c *echo.Context, refreshTokenCookie *http.Cookie) (domain.RefreshResponse, error) {
	found, err := s.store.GetRefreshTokenByHash(c.Request().Context(), refreshTokenCookie.Value)
	if err != nil {
		return domain.RefreshResponse{}, domain.ErrRefreshTokenNotFound
	}

	if found.RevokedAt.Valid || found.ExpiresAt.Time.After(time.Now()) {
		return domain.RefreshResponse{}, domain.ErrRefreshTokenExpired
	}

	_, err = s.store.SetRefreshTokenRevoked(c.Request().Context(), store.SetRefreshTokenRevokedParams{
		TokenHash: refreshTokenCookie.Value,
		RevokedAt: pgtype.Timestamptz{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	})
	if err != nil {
		return domain.RefreshResponse{}, err
	}

	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		return domain.RefreshResponse{}, err
	}

	cookie := new(http.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = newRefreshToken
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.MaxAge = 7 * 24 * 60 * 60

	_, err = s.store.CreateRefreshToken(c.Request().Context(), store.CreateRefreshTokenParams{
		UserID:    found.UserID,
		TokenHash: newRefreshToken,
	})
	if err != nil {
		return domain.RefreshResponse{}, err
	}

	user, err := s.store.GetUserByUUID(c.Request().Context(), found.UserID)
	if err != nil {
		return domain.RefreshResponse{}, err
	}

	token, err := auth.GenerateToken(user.Uuid.String(), string(user.Role))
	if err != nil {
		return domain.RefreshResponse{}, err
	}

	c.SetCookie(cookie)

	return domain.RefreshResponse{
		UUID:     user.Uuid.String(),
		Username: user.Username,
		Role:     string(user.Role),
		JWT:      token,
	}, nil
}

func (s *AuthService) Delete(ctx context.Context, uuid pgtype.UUID) error {
	_, err := s.store.GetUserByUUID(ctx, uuid)
	if err != nil {
		return domain.ErrUsernameDoesNotExists
	}

	_, err = s.store.DeleteUser(ctx, uuid)

	return err
}

func (s *AuthService) PatchProfile(ctx context.Context, uuid pgtype.UUID, newUsername string) error {
	found, err := s.store.GetUserByUUID(ctx, uuid)
	if err != nil {
		return domain.ErrUsernameDoesNotExists
	}

	_, err = s.store.GetUserByUsername(ctx, newUsername)
	if err == nil {
		return domain.ErrUsernameAlreadyExists
	}

	_, err = s.store.UpdateUser(ctx, store.UpdateUserParams{
		Uuid:     uuid,
		Username: newUsername,
		Password: found.Password,
		Role:     found.Role,
		Language: found.Language,
	})

	return err
}

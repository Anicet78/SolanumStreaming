package service

import (
	"context"
	"errors"

	"github.com/Anicet78/SolanumStreaming/auth/internal/domain"
	"github.com/Anicet78/SolanumStreaming/auth/internal/store"
	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	store *store.Queries
}

func NewAuthService(store *store.Queries) *AuthService {
	return &AuthService{store: store}
}

func (s *AuthService) Register(ctx context.Context, username string, password string) (domain.UserResponse, error) {
	exists, err := s.store.GetUserByUsername(ctx, username)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return domain.UserResponse{}, err
		}
	} else if exists.Uuid.Valid {
		return domain.UserResponse{}, domain.ErrUsernameAlreadyExists
	}

	created, err := s.store.CreateUser(ctx, store.CreateUserParams{
		Username: username,
		Password: password,
	})

	if err != nil {
		return domain.UserResponse{}, err
	}

	return domain.UserResponse{
		UUID:     created.Uuid.String(),
		Username: created.Username,
	}, nil
}

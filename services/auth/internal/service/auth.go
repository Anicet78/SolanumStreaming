package service

import (
	"context"
	"errors"
	"shared/auth"

	"github.com/Anicet78/SolanumStreaming/auth/internal/domain"
	"github.com/Anicet78/SolanumStreaming/auth/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthService struct {
	store *store.Queries
}

func NewAuthService(store *store.Queries) *AuthService {
	return &AuthService{store: store}
}

func (s *AuthService) Register(ctx context.Context, username string, password string) (domain.CreateUserResponse, error) {
	exists, err := s.store.GetUserByUsername(ctx, username)
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

	created, err := s.store.CreateUser(ctx, store.CreateUserParams{
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

	return domain.CreateUserResponse{
		UUID:     created.Uuid.String(),
		Username: created.Username,
		Role:     string(created.Role),
		JWT:      token,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, username string, password string) (domain.LoginUserResponse, error) {
	found, err := s.store.GetUserByUsername(ctx, username)
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

	return domain.LoginUserResponse{
		UUID:     found.Uuid.String(),
		Username: found.Username,
		Role:     string(found.Role),
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

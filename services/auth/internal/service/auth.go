package service

import "github.com/Anicet78/SolanumStreaming/auth/internal/store"

type AuthService struct {
	store *store.Queries
}

func NewAuthService(store *store.Queries) *AuthService {

	return &AuthService{store: store}
}

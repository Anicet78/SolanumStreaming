package service

import (
	"errors"
	"os"
	"time"

	"github.com/Anicet78/SolanumStreaming/auth/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

func generateToken(uuid string, role string) (string, error) {
	payload := domain.JWTPayload{
		UUID: uuid,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func verifyToken(tokenString string) (payload *domain.JWTPayload, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTPayload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*domain.JWTPayload)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

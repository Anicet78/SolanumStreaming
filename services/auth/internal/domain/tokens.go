package domain

import "github.com/golang-jwt/jwt/v5"

type JWTPayload struct {
	UUID string `json:"uuid"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

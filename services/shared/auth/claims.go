package auth

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UUID string `json:"uuid"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

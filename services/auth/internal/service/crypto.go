package service

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/alexedwards/argon2id"
)

func hashPassword(password string) (encodedHash string, err error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func passwordMatch(rawPassword string, hashedPassword string) (match bool, err error) {
	return argon2id.ComparePasswordAndHash(rawPassword, hashedPassword)
}

func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 40)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

package service

import "github.com/alexedwards/argon2id"

func hashPassword(password string) (encodedHash string, err error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func passwordMatch(rawPassword string, hashedPassword string) (match bool, err error) {
	return argon2id.ComparePasswordAndHash(rawPassword, hashedPassword)
}

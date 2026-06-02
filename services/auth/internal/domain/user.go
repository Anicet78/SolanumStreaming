package domain

import "errors"

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=4"`
}

type CreateUserResponse struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Role     string `json:"role"`
	JWT      string `json:"jwt"`
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=4"`
}

type LoginUserResponse struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Role     string `json:"role"`
	JWT      string `json:"jwt"`
}

type PatchProfileRequest struct {
	NewUsername string `json:"new_username"`
}

var ErrUsernameAlreadyExists = errors.New("Username already exists")
var ErrUsernameDoesNotExists = errors.New("Username doesn't exists")
var ErrPasswordDoesNotMatch = errors.New("Incorrect password")

package domain

import "errors"

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=4"`
}

type UserResponse struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
}

var ErrUsernameAlreadyExists = errors.New("username already exists")

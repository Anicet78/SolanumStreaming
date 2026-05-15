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

var ErrUsernameAlreadyExists = errors.New("username already exists")
var ErrUsernameDoesNotExists = errors.New("username doesn't exists")
var ErrPasswordDoesNotMatch = errors.New("incorrect password")

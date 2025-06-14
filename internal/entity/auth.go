package entity

import "errors"

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

package jwt

import "errors"

var (
	ErrValidateToken = errors.New("validate token")
	ErrGetClaims     = errors.New("get claims from token")
	ErrGetUserID     = errors.New("get user ID from token")
)

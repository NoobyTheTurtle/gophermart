package middleware

import (
	"github.com/NoobyTheTurtle/gophermart/internal/usecase/auth"
)

type AuthUseCase interface {
	ValidateToken(token string) (int, error)
}

var _ AuthUseCase = (*auth.AuthUseCaseImpl)(nil)

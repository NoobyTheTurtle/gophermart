package middleware

import (
	"github.com/NoobyTheTurtle/gophermart/internal/usecase/auth"
	"github.com/NoobyTheTurtle/gophermart/pkg/logger"
)

type Logger interface {
	Info(message string, args ...any)
	Error(message string, args ...any)
}

var _ Logger = (*logger.ZapLogger)(nil)

type AuthUseCase interface {
	ValidateToken(token string) (int, error)
}

var _ AuthUseCase = (*auth.AuthUseCaseImpl)(nil)

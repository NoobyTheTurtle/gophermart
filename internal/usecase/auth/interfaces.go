package auth

import (
	"context"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByLogin(ctx context.Context, login string) (*entity.User, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
}

type TokenService interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(tokenString string) (int, error)
}

type AuthUseCase interface {
	Register(ctx context.Context, auth *entity.Auth) (token string, err error)
	Login(ctx context.Context, auth *entity.Auth) (token string, err error)
	ValidateToken(token string) (int, error)
}

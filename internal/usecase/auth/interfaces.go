package auth

import (
	"context"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
	"github.com/NoobyTheTurtle/gophermart/internal/repo/postgres/user"
	"github.com/NoobyTheTurtle/gophermart/pkg/jwt"
	"github.com/NoobyTheTurtle/gophermart/pkg/password"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByLogin(ctx context.Context, login string) (*entity.User, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
}

var _ UserRepo = (*user.UserPostgresRepo)(nil)

type TokenService interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(tokenString string) (int, error)
}

var _ TokenService = (*jwt.JWTService)(nil)

type PasswordService interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

var _ PasswordService = (*password.PasswordService)(nil)

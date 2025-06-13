package auth

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
)

type AuthUseCaseImpl struct {
	userRepo     UserRepo
	tokenService TokenService
}

func New(userRepo UserRepo, tokenService TokenService) *AuthUseCaseImpl {
	return &AuthUseCaseImpl{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// Register creates new user and returns JWT token
func (a *AuthUseCaseImpl) Register(ctx context.Context, auth *entity.Auth) (string, error) {
	existingUser, err := a.userRepo.GetUserByLogin(ctx, auth.Login)
	if err == nil && existingUser != nil {
		return "", ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(auth.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("usecase.Register: failed to hash password: %w", err)
	}

	user := &entity.User{
		Login:        auth.Login,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	if err = a.userRepo.CreateUser(ctx, user); err != nil {
		return "", fmt.Errorf("usecase.Register: failed to create user: %w", err)
	}

	token, err := a.tokenService.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("usecase.Register: failed to generate token: %w", err)
	}

	return token, nil
}

// Login authenticates user and returns JWT token
func (a *AuthUseCaseImpl) Login(ctx context.Context, auth *entity.Auth) (string, error) {
	user, err := a.userRepo.GetUserByLogin(ctx, auth.Login)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(auth.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := a.tokenService.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("usecase.Login: failed to generate token: %w", err)
	}

	return token, nil
}

// ValidateToken validates JWT token and returns user ID
func (a *AuthUseCaseImpl) ValidateToken(tokenString string) (int, error) {
	return a.tokenService.ValidateToken(tokenString)
}

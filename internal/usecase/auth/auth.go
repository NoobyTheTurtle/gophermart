package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
)

type AuthUseCaseImpl struct {
	userRepo        UserRepo
	tokenService    TokenService
	passwordService PasswordService
}

func New(userRepo UserRepo, tokenService TokenService, passwordService PasswordService) *AuthUseCaseImpl {
	return &AuthUseCaseImpl{
		userRepo:        userRepo,
		tokenService:    tokenService,
		passwordService: passwordService,
	}
}

func (a *AuthUseCaseImpl) Register(ctx context.Context, auth *entity.Auth) (string, error) {
	existingUser, err := a.userRepo.GetUserByLogin(ctx, auth.Login)
	if err == nil && existingUser != nil {
		return "", entity.ErrUserAlreadyExists
	}

	hashedPassword, err := a.passwordService.Hash(auth.Password)
	if err != nil {
		return "", fmt.Errorf("authUseCase.Register: failed to hash password: %w", err)
	}

	user := &entity.User{
		Login:        auth.Login,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}

	if err = a.userRepo.CreateUser(ctx, user); err != nil {
		return "", fmt.Errorf("authUseCase.Register: failed to create user: %w", err)
	}

	token, err := a.tokenService.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("authUseCase.Register: failed to generate token: %w", err)
	}

	return token, nil
}

func (a *AuthUseCaseImpl) Login(ctx context.Context, auth *entity.Auth) (string, error) {
	user, err := a.userRepo.GetUserByLogin(ctx, auth.Login)
	if err != nil {
		return "", entity.ErrInvalidCredentials
	}

	if err = a.passwordService.Compare(user.PasswordHash, auth.Password); err != nil {
		return "", entity.ErrInvalidCredentials
	}

	token, err := a.tokenService.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("authUseCase.Login: failed to generate token: %w", err)
	}

	return token, nil
}

func (a *AuthUseCaseImpl) ValidateToken(tokenString string) (int, error) {
	return a.tokenService.ValidateToken(tokenString)
}

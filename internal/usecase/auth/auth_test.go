package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthUseCaseImpl_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := NewMockUserRepo(ctrl)
	mockTokenService := NewMockTokenService(ctrl)
	mockPasswordService := NewMockPasswordService(ctrl)

	uc := New(mockUserRepo, mockTokenService, mockPasswordService)

	tests := []struct {
		name       string
		auth       *entity.Auth
		setupMocks func()
		want       string
		wantErr    bool
		errType    error
	}{
		{
			name: "successful registration",
			auth: &entity.Auth{
				Login:    "testuser",
				Password: "password123",
			},
			setupMocks: func() {
				mockUserRepo.EXPECT().GetUserByLogin(gomock.Any(), "testuser").Return(nil, entity.ErrUserNotFound)
				mockPasswordService.EXPECT().Hash("password123").Return("hashed_password", nil)
				mockUserRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, user *entity.User) error {
						user.ID = 1
						return nil
					})
				mockTokenService.EXPECT().GenerateToken(1).Return("test_token", nil)
			},
			want:    "test_token",
			wantErr: false,
		},
		{
			name: "user already exists",
			auth: &entity.Auth{
				Login:    "testuser",
				Password: "password123",
			},
			setupMocks: func() {
				existingUser := &entity.User{
					ID:           1,
					Login:        "testuser",
					PasswordHash: "existing_hash",
				}
				mockUserRepo.EXPECT().GetUserByLogin(gomock.Any(), "testuser").Return(existingUser, nil)
			},
			want:    "",
			wantErr: true,
			errType: entity.ErrUserAlreadyExists,
		},
		{
			name: "password hash error",
			auth: &entity.Auth{
				Login:    "testuser",
				Password: "password123",
			},
			setupMocks: func() {
				mockUserRepo.EXPECT().GetUserByLogin(gomock.Any(), "testuser").Return(nil, entity.ErrUserNotFound)
				mockPasswordService.EXPECT().Hash("password123").Return("", errors.New("hash error"))
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "create user error",
			auth: &entity.Auth{
				Login:    "testuser",
				Password: "password123",
			},
			setupMocks: func() {
				mockUserRepo.EXPECT().GetUserByLogin(gomock.Any(), "testuser").Return(nil, entity.ErrUserNotFound)
				mockPasswordService.EXPECT().Hash("password123").Return("hashed_password", nil)
				mockUserRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(errors.New("create error"))
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "token generation error",
			auth: &entity.Auth{
				Login:    "testuser",
				Password: "password123",
			},
			setupMocks: func() {
				mockUserRepo.EXPECT().GetUserByLogin(gomock.Any(), "testuser").Return(nil, entity.ErrUserNotFound)
				mockPasswordService.EXPECT().Hash("password123").Return("hashed_password", nil)
				mockUserRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, user *entity.User) error {
						user.ID = 1
						return nil
					})
				mockTokenService.EXPECT().GenerateToken(1).Return("", errors.New("token error"))
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			got, err := uc.Register(context.Background(), tt.auth)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestAuthUseCaseImpl_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := NewMockUserRepo(ctrl)
	mockTokenService := NewMockTokenService(ctrl)
	mockPasswordService := NewMockPasswordService(ctrl)

	uc := New(mockUserRepo, mockTokenService, mockPasswordService)

	tests := []struct {
		name       string
		auth       *entity.Auth
		setupMocks func()
		want       string
		wantErr    bool
		errType    error
	}{
		{
			name: "successful login",
			auth: &entity.Auth{
				Login:    "testuser",
				Password: "password123",
			},
			setupMocks: func() {
				user := &entity.User{
					ID:           1,
					Login:        "testuser",
					PasswordHash: "hashed_password",
					CreatedAt:    time.Now(),
				}
				mockUserRepo.EXPECT().GetUserByLogin(gomock.Any(), "testuser").Return(user, nil)
				mockPasswordService.EXPECT().Compare("hashed_password", "password123").Return(nil)
				mockTokenService.EXPECT().GenerateToken(1).Return("test_token", nil)
			},
			want:    "test_token",
			wantErr: false,
		},
		{
			name: "user not found",
			auth: &entity.Auth{
				Login:    "testuser",
				Password: "password123",
			},
			setupMocks: func() {
				mockUserRepo.EXPECT().GetUserByLogin(gomock.Any(), "testuser").Return(nil, entity.ErrUserNotFound)
			},
			want:    "",
			wantErr: true,
			errType: entity.ErrInvalidCredentials,
		},
		{
			name: "invalid password",
			auth: &entity.Auth{
				Login:    "testuser",
				Password: "wrongpassword",
			},
			setupMocks: func() {
				user := &entity.User{
					ID:           1,
					Login:        "testuser",
					PasswordHash: "hashed_password",
					CreatedAt:    time.Now(),
				}
				mockUserRepo.EXPECT().GetUserByLogin(gomock.Any(), "testuser").Return(user, nil)
				mockPasswordService.EXPECT().Compare("hashed_password", "wrongpassword").Return(errors.New("invalid password"))
			},
			want:    "",
			wantErr: true,
			errType: entity.ErrInvalidCredentials,
		},
		{
			name: "token generation error",
			auth: &entity.Auth{
				Login:    "testuser",
				Password: "password123",
			},
			setupMocks: func() {
				user := &entity.User{
					ID:           1,
					Login:        "testuser",
					PasswordHash: "hashed_password",
					CreatedAt:    time.Now(),
				}
				mockUserRepo.EXPECT().GetUserByLogin(gomock.Any(), "testuser").Return(user, nil)
				mockPasswordService.EXPECT().Compare("hashed_password", "password123").Return(nil)
				mockTokenService.EXPECT().GenerateToken(1).Return("", errors.New("token error"))
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			got, err := uc.Login(context.Background(), tt.auth)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestAuthUseCaseImpl_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := NewMockUserRepo(ctrl)
	mockTokenService := NewMockTokenService(ctrl)
	mockPasswordService := NewMockPasswordService(ctrl)

	uc := New(mockUserRepo, mockTokenService, mockPasswordService)

	tests := []struct {
		name        string
		tokenString string
		setupMocks  func()
		want        int
		wantErr     bool
	}{
		{
			name:        "successful token validation",
			tokenString: "valid_token",
			setupMocks: func() {
				mockTokenService.EXPECT().ValidateToken("valid_token").Return(1, nil)
			},
			want:    1,
			wantErr: false,
		},
		{
			name:        "invalid token",
			tokenString: "invalid_token",
			setupMocks: func() {
				mockTokenService.EXPECT().ValidateToken("invalid_token").Return(0, errors.New("invalid token"))
			},
			want:    0,
			wantErr: true,
		},
		{
			name:        "empty token",
			tokenString: "",
			setupMocks: func() {
				mockTokenService.EXPECT().ValidateToken("").Return(0, errors.New("empty token"))
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			got, err := uc.ValidateToken(tt.tokenString)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

package balance

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBalanceUseCaseImpl_GetUserBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBalanceRepo := NewMockBalanceRepo(ctrl)
	uc := New(mockBalanceRepo)

	tests := []struct {
		name       string
		userID     int
		setupMocks func()
		want       *entity.UserBalance
		wantErr    bool
	}{
		{
			name:   "successful get existing balance",
			userID: 1,
			setupMocks: func() {
				balance := &entity.UserBalance{
					UserID:    1,
					Current:   100.50,
					Withdrawn: 25.00,
					UpdatedAt: time.Now(),
				}
				mockBalanceRepo.EXPECT().GetBalanceByUserID(gomock.Any(), 1).Return(balance, nil)
			},
			want: &entity.UserBalance{
				UserID:    1,
				Current:   100.50,
				Withdrawn: 25.00,
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name:   "balance not found - create new balance",
			userID: 2,
			setupMocks: func() {
				mockBalanceRepo.EXPECT().GetBalanceByUserID(gomock.Any(), 2).Return(nil, entity.ErrBalanceNotFound)
				mockBalanceRepo.EXPECT().CreateBalance(gomock.Any(), 2).Return(nil)
			},
			want: &entity.UserBalance{
				UserID:    2,
				Current:   0,
				Withdrawn: 0,
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name:   "balance not found - create balance error",
			userID: 3,
			setupMocks: func() {
				mockBalanceRepo.EXPECT().GetBalanceByUserID(gomock.Any(), 3).Return(nil, entity.ErrBalanceNotFound)
				mockBalanceRepo.EXPECT().CreateBalance(gomock.Any(), 3).Return(errors.New("create error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "repository error",
			userID: 4,
			setupMocks: func() {
				mockBalanceRepo.EXPECT().GetBalanceByUserID(gomock.Any(), 4).Return(nil, errors.New("repo error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			got, err := uc.GetUserBalance(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.UserID, got.UserID)
				assert.Equal(t, tt.want.Current, got.Current)
				assert.Equal(t, tt.want.Withdrawn, got.Withdrawn)
				assert.Equal(t, tt.want.UpdatedAt.Unix(), got.UpdatedAt.Unix())
			}
		})
	}
}

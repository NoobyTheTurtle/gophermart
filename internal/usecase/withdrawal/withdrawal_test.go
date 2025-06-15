package withdrawal

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestWithdrawalUseCaseImpl_ProcessWithdrawal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawalRepo := NewMockWithdrawalRepo(ctrl)
	mockBalanceRepo := NewMockBalanceRepo(ctrl)
	mockLuhnService := NewMockLuhnService(ctrl)

	uc := New(mockWithdrawalRepo, mockBalanceRepo, mockLuhnService)

	tests := []struct {
		name        string
		userID      int
		orderNumber string
		sum         float64
		setupMocks  func()
		wantErr     bool
		errType     error
	}{
		{
			name:        "successful withdrawal",
			userID:      1,
			orderNumber: "1234567890",
			sum:         50.0,
			setupMocks: func() {
				mockLuhnService.EXPECT().IsValid("1234567890").Return(true)
				balance := &entity.UserBalance{
					UserID:    1,
					Current:   100.0,
					Withdrawn: 25.0,
					UpdatedAt: time.Now(),
				}
				mockBalanceRepo.EXPECT().GetBalanceByUserID(gomock.Any(), 1).Return(balance, nil)
				mockWithdrawalRepo.EXPECT().CreateWithdrawal(gomock.Any(), gomock.Any()).Return(nil)
				mockBalanceRepo.EXPECT().UpdateBalance(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, b *entity.UserBalance) error {
						assert.Equal(t, 50.0, b.Current)
						assert.Equal(t, 75.0, b.Withdrawn)
						return nil
					})
			},
			wantErr: false,
		},
		{
			name:        "invalid order number",
			userID:      1,
			orderNumber: "invalid",
			sum:         50.0,
			setupMocks: func() {
				mockLuhnService.EXPECT().IsValid("invalid").Return(false)
			},
			wantErr: true,
			errType: entity.ErrInvalidOrderNumber,
		},
		{
			name:        "insufficient funds",
			userID:      1,
			orderNumber: "1234567890",
			sum:         150.0,
			setupMocks: func() {
				mockLuhnService.EXPECT().IsValid("1234567890").Return(true)
				balance := &entity.UserBalance{
					UserID:    1,
					Current:   100.0,
					Withdrawn: 25.0,
					UpdatedAt: time.Now(),
				}
				mockBalanceRepo.EXPECT().GetBalanceByUserID(gomock.Any(), 1).Return(balance, nil)
			},
			wantErr: true,
			errType: entity.ErrInsufficientFunds,
		},
		{
			name:        "get balance error",
			userID:      1,
			orderNumber: "1234567890",
			sum:         50.0,
			setupMocks: func() {
				mockLuhnService.EXPECT().IsValid("1234567890").Return(true)
				mockBalanceRepo.EXPECT().GetBalanceByUserID(gomock.Any(), 1).Return(nil, errors.New("balance error"))
			},
			wantErr: true,
		},
		{
			name:        "create withdrawal error",
			userID:      1,
			orderNumber: "1234567890",
			sum:         50.0,
			setupMocks: func() {
				mockLuhnService.EXPECT().IsValid("1234567890").Return(true)
				balance := &entity.UserBalance{
					UserID:    1,
					Current:   100.0,
					Withdrawn: 25.0,
					UpdatedAt: time.Now(),
				}
				mockBalanceRepo.EXPECT().GetBalanceByUserID(gomock.Any(), 1).Return(balance, nil)
				mockWithdrawalRepo.EXPECT().CreateWithdrawal(gomock.Any(), gomock.Any()).Return(errors.New("create error"))
			},
			wantErr: true,
		},
		{
			name:        "update balance error",
			userID:      1,
			orderNumber: "1234567890",
			sum:         50.0,
			setupMocks: func() {
				mockLuhnService.EXPECT().IsValid("1234567890").Return(true)
				balance := &entity.UserBalance{
					UserID:    1,
					Current:   100.0,
					Withdrawn: 25.0,
					UpdatedAt: time.Now(),
				}
				mockBalanceRepo.EXPECT().GetBalanceByUserID(gomock.Any(), 1).Return(balance, nil)
				mockWithdrawalRepo.EXPECT().CreateWithdrawal(gomock.Any(), gomock.Any()).Return(nil)
				mockBalanceRepo.EXPECT().UpdateBalance(gomock.Any(), gomock.Any()).Return(errors.New("update error"))
			},
			wantErr: true,
		},
		{
			name:        "zero withdrawal amount",
			userID:      1,
			orderNumber: "1234567890",
			sum:         0.0,
			setupMocks: func() {
				mockLuhnService.EXPECT().IsValid("1234567890").Return(true)
				balance := &entity.UserBalance{
					UserID:    1,
					Current:   100.0,
					Withdrawn: 25.0,
					UpdatedAt: time.Now(),
				}
				mockBalanceRepo.EXPECT().GetBalanceByUserID(gomock.Any(), 1).Return(balance, nil)
				mockWithdrawalRepo.EXPECT().CreateWithdrawal(gomock.Any(), gomock.Any()).Return(nil)
				mockBalanceRepo.EXPECT().UpdateBalance(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, b *entity.UserBalance) error {
						assert.Equal(t, 100.0, b.Current)
						assert.Equal(t, 25.0, b.Withdrawn)
						return nil
					})
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := uc.ProcessWithdrawal(context.Background(), tt.userID, tt.orderNumber, tt.sum)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWithdrawalUseCaseImpl_GetUserWithdrawals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawalRepo := NewMockWithdrawalRepo(ctrl)
	mockBalanceRepo := NewMockBalanceRepo(ctrl)
	mockLuhnService := NewMockLuhnService(ctrl)

	uc := New(mockWithdrawalRepo, mockBalanceRepo, mockLuhnService)

	tests := []struct {
		name       string
		userID     int
		setupMocks func()
		want       []*entity.Withdrawal
		wantErr    bool
	}{
		{
			name:   "successful get user withdrawals",
			userID: 1,
			setupMocks: func() {
				withdrawals := []*entity.Withdrawal{
					{
						ID:          1,
						UserID:      1,
						OrderNumber: "1234567890",
						Sum:         50.0,
						ProcessedAt: time.Now(),
					},
					{
						ID:          2,
						UserID:      1,
						OrderNumber: "0987654321",
						Sum:         25.0,
						ProcessedAt: time.Now(),
					},
				}
				mockWithdrawalRepo.EXPECT().GetWithdrawalsByUserID(gomock.Any(), 1).Return(withdrawals, nil)
			},
			want: []*entity.Withdrawal{
				{
					ID:          1,
					UserID:      1,
					OrderNumber: "1234567890",
					Sum:         50.0,
					ProcessedAt: time.Now(),
				},
				{
					ID:          2,
					UserID:      1,
					OrderNumber: "0987654321",
					Sum:         25.0,
					ProcessedAt: time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name:   "empty withdrawals list",
			userID: 2,
			setupMocks: func() {
				mockWithdrawalRepo.EXPECT().GetWithdrawalsByUserID(gomock.Any(), 2).Return([]*entity.Withdrawal{}, nil)
			},
			want:    []*entity.Withdrawal{},
			wantErr: false,
		},
		{
			name:   "repository error",
			userID: 3,
			setupMocks: func() {
				mockWithdrawalRepo.EXPECT().GetWithdrawalsByUserID(gomock.Any(), 3).Return(nil, errors.New("repo error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			got, err := uc.GetUserWithdrawals(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.want), len(got))
				if len(tt.want) > 0 {
					for i, withdrawal := range got {
						assert.Equal(t, tt.want[i].ID, withdrawal.ID)
						assert.Equal(t, tt.want[i].UserID, withdrawal.UserID)
						assert.Equal(t, tt.want[i].OrderNumber, withdrawal.OrderNumber)
						assert.Equal(t, tt.want[i].Sum, withdrawal.Sum)
					}
				}
			}
		})
	}
}

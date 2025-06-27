package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOrderUseCaseImpl_UploadOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := NewMockOrderRepo(ctrl)
	mockBalanceRepo := NewMockBalanceRepo(ctrl)
	mockLuhnService := NewMockLuhnService(ctrl)

	uc := New(mockOrderRepo, mockBalanceRepo, mockLuhnService)

	tests := []struct {
		name        string
		userID      int
		orderNumber string
		setupMocks  func()
		wantErr     bool
		errType     error
	}{
		{
			name:        "successful order upload",
			userID:      1,
			orderNumber: "1234567890",
			setupMocks: func() {
				mockLuhnService.EXPECT().IsValid("1234567890").Return(true)
				mockOrderRepo.EXPECT().GetOrderByNumber(gomock.Any(), "1234567890").Return(nil, entity.ErrOrderNotFound)
				mockOrderRepo.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(nil)
				mockBalanceRepo.EXPECT().CreateBalance(gomock.Any(), 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "invalid order number",
			userID:      1,
			orderNumber: "invalid",
			setupMocks: func() {
				mockLuhnService.EXPECT().IsValid("invalid").Return(false)
			},
			wantErr: true,
			errType: entity.ErrInvalidOrderNumber,
		},
		{
			name:        "order already exists for same user",
			userID:      1,
			orderNumber: "1234567890",
			setupMocks: func() {
				mockLuhnService.EXPECT().IsValid("1234567890").Return(true)
				existingOrder := &entity.Order{
					ID:     1,
					UserID: 1,
					Number: "1234567890",
					Status: entity.OrderStatusNew,
				}
				mockOrderRepo.EXPECT().GetOrderByNumber(gomock.Any(), "1234567890").Return(existingOrder, nil)
			},
			wantErr: true,
			errType: entity.ErrOrderAlreadyExists,
		},
		{
			name:        "order already exists for another user",
			userID:      1,
			orderNumber: "1234567890",
			setupMocks: func() {
				mockLuhnService.EXPECT().IsValid("1234567890").Return(true)
				existingOrder := &entity.Order{
					ID:     1,
					UserID: 2,
					Number: "1234567890",
					Status: entity.OrderStatusNew,
				}
				mockOrderRepo.EXPECT().GetOrderByNumber(gomock.Any(), "1234567890").Return(existingOrder, nil)
			},
			wantErr: true,
			errType: entity.ErrOrderAlreadyExistsByAnotherUser,
		},
		{
			name:        "repository error checking order existence",
			userID:      1,
			orderNumber: "1234567890",
			setupMocks: func() {
				mockLuhnService.EXPECT().IsValid("1234567890").Return(true)
				mockOrderRepo.EXPECT().GetOrderByNumber(gomock.Any(), "1234567890").Return(nil, errors.New("repo error"))
			},
			wantErr: true,
		},
		{
			name:        "create order error",
			userID:      1,
			orderNumber: "1234567890",
			setupMocks: func() {
				mockLuhnService.EXPECT().IsValid("1234567890").Return(true)
				mockOrderRepo.EXPECT().GetOrderByNumber(gomock.Any(), "1234567890").Return(nil, entity.ErrOrderNotFound)
				mockOrderRepo.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(errors.New("create error"))
			},
			wantErr: true,
		},
		{
			name:        "create balance error",
			userID:      1,
			orderNumber: "1234567890",
			setupMocks: func() {
				mockLuhnService.EXPECT().IsValid("1234567890").Return(true)
				mockOrderRepo.EXPECT().GetOrderByNumber(gomock.Any(), "1234567890").Return(nil, entity.ErrOrderNotFound)
				mockOrderRepo.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(nil)
				mockBalanceRepo.EXPECT().CreateBalance(gomock.Any(), 1).Return(errors.New("balance error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := uc.UploadOrder(context.Background(), tt.userID, tt.orderNumber)

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

func TestOrderUseCaseImpl_GetUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := NewMockOrderRepo(ctrl)
	mockBalanceRepo := NewMockBalanceRepo(ctrl)
	mockLuhnService := NewMockLuhnService(ctrl)

	uc := New(mockOrderRepo, mockBalanceRepo, mockLuhnService)

	tests := []struct {
		name       string
		userID     int
		setupMocks func()
		want       []*entity.Order
		wantErr    bool
	}{
		{
			name:   "successful get user orders",
			userID: 1,
			setupMocks: func() {
				orders := []*entity.Order{
					{
						ID:         1,
						UserID:     1,
						Number:     "1234567890",
						Status:     entity.OrderStatusProcessed,
						Accrual:    100.50,
						UploadedAt: time.Now(),
					},
					{
						ID:         2,
						UserID:     1,
						Number:     "0987654321",
						Status:     entity.OrderStatusNew,
						Accrual:    0,
						UploadedAt: time.Now(),
					},
				}
				mockOrderRepo.EXPECT().GetOrdersByUserID(gomock.Any(), 1).Return(orders, nil)
			},
			want: []*entity.Order{
				{
					ID:         1,
					UserID:     1,
					Number:     "1234567890",
					Status:     entity.OrderStatusProcessed,
					Accrual:    100.50,
					UploadedAt: time.Now(),
				},
				{
					ID:         2,
					UserID:     1,
					Number:     "0987654321",
					Status:     entity.OrderStatusNew,
					Accrual:    0,
					UploadedAt: time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name:   "empty orders list",
			userID: 2,
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrdersByUserID(gomock.Any(), 2).Return([]*entity.Order{}, nil)
			},
			want:    []*entity.Order{},
			wantErr: false,
		},
		{
			name:   "repository error",
			userID: 3,
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrdersByUserID(gomock.Any(), 3).Return(nil, errors.New("repo error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			got, err := uc.GetUserOrders(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.want), len(got))
				if len(tt.want) > 0 {
					for i, order := range got {
						assert.Equal(t, tt.want[i].ID, order.ID)
						assert.Equal(t, tt.want[i].UserID, order.UserID)
						assert.Equal(t, tt.want[i].Number, order.Number)
						assert.Equal(t, tt.want[i].Status, order.Status)
						assert.Equal(t, tt.want[i].Accrual, order.Accrual)
					}
				}
			}
		})
	}
}

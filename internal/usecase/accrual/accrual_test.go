package accrual

import (
	"context"
	"errors"
	"testing"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAccrualUseCaseImpl_ProcessOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccrualAPI := NewMockAccrualAPI(ctrl)
	mockOrderRepo := NewMockOrderRepo(ctrl)
	mockBalanceRepo := NewMockBalanceRepo(ctrl)

	uc := New(mockAccrualAPI, mockOrderRepo, mockBalanceRepo)

	tests := []struct {
		name        string
		orderNumber string
		setupMocks  func()
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "successful processing with accrual",
			orderNumber: "123456789",
			setupMocks: func() {
				order := &entity.Order{
					ID:     1,
					UserID: 1,
					Number: "123456789",
					Status: entity.OrderStatusNew,
				}
				accrual := &entity.Accrual{
					Order:   "123456789",
					Status:  entity.AccrualStatusProcessed,
					Accrual: 100.50,
				}

				mockOrderRepo.EXPECT().GetOrderByNumber(gomock.Any(), "123456789").Return(order, nil)
				mockAccrualAPI.EXPECT().GetOrderAccrual(gomock.Any(), "123456789").Return(accrual, nil)
				mockOrderRepo.EXPECT().UpdateOrderStatus(gomock.Any(), "123456789", entity.OrderStatusProcessed, 100.50).Return(nil)
				mockBalanceRepo.EXPECT().AddAccrual(gomock.Any(), 1, 100.50).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "successful processing without accrual",
			orderNumber: "123456789",
			setupMocks: func() {
				order := &entity.Order{
					ID:     1,
					UserID: 1,
					Number: "123456789",
					Status: entity.OrderStatusNew,
				}
				accrual := &entity.Accrual{
					Order:   "123456789",
					Status:  entity.AccrualStatusProcessed,
					Accrual: 0,
				}

				mockOrderRepo.EXPECT().GetOrderByNumber(gomock.Any(), "123456789").Return(order, nil)
				mockAccrualAPI.EXPECT().GetOrderAccrual(gomock.Any(), "123456789").Return(accrual, nil)
				mockOrderRepo.EXPECT().UpdateOrderStatus(gomock.Any(), "123456789", entity.OrderStatusProcessed, 0.0).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "processing order with invalid status",
			orderNumber: "123456789",
			setupMocks: func() {
				order := &entity.Order{
					ID:     1,
					UserID: 1,
					Number: "123456789",
					Status: entity.OrderStatusNew,
				}
				accrual := &entity.Accrual{
					Order:   "123456789",
					Status:  entity.AccrualStatusInvalid,
					Accrual: 0,
				}

				mockOrderRepo.EXPECT().GetOrderByNumber(gomock.Any(), "123456789").Return(order, nil)
				mockAccrualAPI.EXPECT().GetOrderAccrual(gomock.Any(), "123456789").Return(accrual, nil)
				mockOrderRepo.EXPECT().UpdateOrderStatus(gomock.Any(), "123456789", entity.OrderStatusInvalid, 0.0).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "processing order with registered status",
			orderNumber: "123456789",
			setupMocks: func() {
				order := &entity.Order{
					ID:     1,
					UserID: 1,
					Number: "123456789",
					Status: entity.OrderStatusNew,
				}
				accrual := &entity.Accrual{
					Order:   "123456789",
					Status:  entity.AccrualStatusRegistered,
					Accrual: 0,
				}

				mockOrderRepo.EXPECT().GetOrderByNumber(gomock.Any(), "123456789").Return(order, nil)
				mockAccrualAPI.EXPECT().GetOrderAccrual(gomock.Any(), "123456789").Return(accrual, nil)
				mockOrderRepo.EXPECT().UpdateOrderStatus(gomock.Any(), "123456789", entity.OrderStatusProcessing, 0.0).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "order not found",
			orderNumber: "123456789",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByNumber(gomock.Any(), "123456789").Return(nil, entity.ErrOrderNotFound)
			},
			wantErr: true,
		},
		{
			name:        "accrual API error",
			orderNumber: "123456789",
			setupMocks: func() {
				order := &entity.Order{
					ID:     1,
					UserID: 1,
					Number: "123456789",
					Status: entity.OrderStatusNew,
				}

				mockOrderRepo.EXPECT().GetOrderByNumber(gomock.Any(), "123456789").Return(order, nil)
				mockAccrualAPI.EXPECT().GetOrderAccrual(gomock.Any(), "123456789").Return(nil, errors.New("api error"))
			},
			wantErr: true,
		},
		{
			name:        "unknown accrual status",
			orderNumber: "123456789",
			setupMocks: func() {
				order := &entity.Order{
					ID:     1,
					UserID: 1,
					Number: "123456789",
					Status: entity.OrderStatusNew,
				}
				accrual := &entity.Accrual{
					Order:   "123456789",
					Status:  "unknown_status",
					Accrual: 0,
				}

				mockOrderRepo.EXPECT().GetOrderByNumber(gomock.Any(), "123456789").Return(order, nil)
				mockAccrualAPI.EXPECT().GetOrderAccrual(gomock.Any(), "123456789").Return(accrual, nil)
			},
			wantErr: true,
			errMsg:  "unexpected accrual status",
		},
		{
			name:        "update order status error",
			orderNumber: "123456789",
			setupMocks: func() {
				order := &entity.Order{
					ID:     1,
					UserID: 1,
					Number: "123456789",
					Status: entity.OrderStatusNew,
				}
				accrual := &entity.Accrual{
					Order:   "123456789",
					Status:  entity.AccrualStatusProcessed,
					Accrual: 100.50,
				}

				mockOrderRepo.EXPECT().GetOrderByNumber(gomock.Any(), "123456789").Return(order, nil)
				mockAccrualAPI.EXPECT().GetOrderAccrual(gomock.Any(), "123456789").Return(accrual, nil)
				mockOrderRepo.EXPECT().UpdateOrderStatus(gomock.Any(), "123456789", entity.OrderStatusProcessed, 100.50).Return(errors.New("update error"))
			},
			wantErr: true,
		},
		{
			name:        "add accrual error",
			orderNumber: "123456789",
			setupMocks: func() {
				order := &entity.Order{
					ID:     1,
					UserID: 1,
					Number: "123456789",
					Status: entity.OrderStatusNew,
				}
				accrual := &entity.Accrual{
					Order:   "123456789",
					Status:  entity.AccrualStatusProcessed,
					Accrual: 100.50,
				}

				mockOrderRepo.EXPECT().GetOrderByNumber(gomock.Any(), "123456789").Return(order, nil)
				mockAccrualAPI.EXPECT().GetOrderAccrual(gomock.Any(), "123456789").Return(accrual, nil)
				mockOrderRepo.EXPECT().UpdateOrderStatus(gomock.Any(), "123456789", entity.OrderStatusProcessed, 100.50).Return(nil)
				mockBalanceRepo.EXPECT().AddAccrual(gomock.Any(), 1, 100.50).Return(errors.New("balance error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := uc.ProcessOrder(context.Background(), tt.orderNumber)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAccrualUseCaseImpl_GetOrdersForProcessing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccrualAPI := NewMockAccrualAPI(ctrl)
	mockOrderRepo := NewMockOrderRepo(ctrl)
	mockBalanceRepo := NewMockBalanceRepo(ctrl)

	uc := New(mockAccrualAPI, mockOrderRepo, mockBalanceRepo)

	tests := []struct {
		name       string
		setupMocks func()
		want       []*entity.Order
		wantErr    bool
	}{
		{
			name: "successful get orders",
			setupMocks: func() {
				orders := []*entity.Order{
					{ID: 1, Number: "123456789", Status: entity.OrderStatusNew},
					{ID: 2, Number: "987654321", Status: entity.OrderStatusProcessing},
				}
				mockOrderRepo.EXPECT().GetOrdersForProcessing(gomock.Any()).Return(orders, nil)
			},
			want: []*entity.Order{
				{ID: 1, Number: "123456789", Status: entity.OrderStatusNew},
				{ID: 2, Number: "987654321", Status: entity.OrderStatusProcessing},
			},
			wantErr: false,
		},
		{
			name: "repository error",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrdersForProcessing(gomock.Any()).Return(nil, errors.New("repo error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty orders list",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrdersForProcessing(gomock.Any()).Return([]*entity.Order{}, nil)
			},
			want:    []*entity.Order{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			got, err := uc.GetOrdersForProcessing(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

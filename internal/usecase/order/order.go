package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
)

type OrderUseCaseImpl struct {
	orderRepo   OrderRepo
	balanceRepo BalanceRepo
	luhnService LuhnService
}

func New(orderRepo OrderRepo, balanceRepo BalanceRepo, luhnService LuhnService) *OrderUseCaseImpl {
	return &OrderUseCaseImpl{
		orderRepo:   orderRepo,
		balanceRepo: balanceRepo,
		luhnService: luhnService,
	}
}

func (uc *OrderUseCaseImpl) UploadOrder(ctx context.Context, userID int, orderNumber string) error {
	if !uc.luhnService.IsValid(orderNumber) {
		return entity.ErrInvalidOrderNumber
	}

	existingOrder, err := uc.orderRepo.GetOrderByNumber(ctx, orderNumber)
	if err != nil && !errors.Is(err, entity.ErrOrderNotFound) {
		return fmt.Errorf("usecase - order - UploadOrder: failed to check order existence: %w", err)
	}

	if existingOrder != nil {
		if existingOrder.UserID == userID {
			return entity.ErrOrderAlreadyExists
		}
		return entity.ErrOrderAlreadyExistsByAnotherUser
	}

	newOrder := &entity.Order{
		UserID:     userID,
		Number:     orderNumber,
		Status:     entity.OrderStatusNew,
		UploadedAt: time.Now(),
	}

	if err := uc.orderRepo.CreateOrder(ctx, newOrder); err != nil {
		return fmt.Errorf("usecase - order - UploadOrder: failed to create order: %w", err)
	}

	if err := uc.balanceRepo.CreateBalance(ctx, userID); err != nil {
		return fmt.Errorf("usecase - order - UploadOrder: failed to create balance: %w", err)
	}

	return nil
}

func (uc *OrderUseCaseImpl) GetUserOrders(ctx context.Context, userID int) ([]*entity.Order, error) {
	orders, err := uc.orderRepo.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("usecase - order - GetUserOrders: failed to get user orders: %w", err)
	}

	return orders, nil
}

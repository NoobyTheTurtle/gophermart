package accrual

import (
	"context"
	"fmt"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
)

type AccrualUseCaseImpl struct {
	accrualAPI  AccrualAPI
	orderRepo   OrderRepo
	balanceRepo BalanceRepo
}

func New(
	accrualAPI AccrualAPI,
	orderRepo OrderRepo,
	balanceRepo BalanceRepo,
) *AccrualUseCaseImpl {
	return &AccrualUseCaseImpl{
		accrualAPI:  accrualAPI,
		orderRepo:   orderRepo,
		balanceRepo: balanceRepo,
	}
}

func (s *AccrualUseCaseImpl) ProcessOrder(ctx context.Context, orderNumber string) error {
	order, err := s.orderRepo.GetOrderByNumber(ctx, orderNumber)
	if err != nil {
		return err
	}

	accrual, err := s.accrualAPI.GetOrderAccrual(ctx, orderNumber)
	if err != nil {
		return err
	}

	var newStatus entity.OrderStatus
	switch accrual.Status {
	case entity.AccrualStatusInvalid:
		newStatus = entity.OrderStatusInvalid
	case entity.AccrualStatusProcessed:
		newStatus = entity.OrderStatusProcessed
	case entity.AccrualStatusProcessing:
		newStatus = entity.OrderStatusProcessing
	case entity.AccrualStatusRegistered:
		newStatus = entity.OrderStatusProcessing
	default:
		return fmt.Errorf("usecase - accrual - ProcessOrder: unexpected accrual status: %s", accrual.Status)
	}

	if err := s.orderRepo.UpdateOrderStatus(ctx, orderNumber, newStatus, accrual.Accrual); err != nil {
		return err
	}

	if newStatus == entity.OrderStatusProcessed && accrual.Accrual > 0 {
		if err := s.balanceRepo.AddAccrual(ctx, order.UserID, accrual.Accrual); err != nil {
			return err
		}
	}

	return nil
}

func (s *AccrualUseCaseImpl) GetOrdersForProcessing(ctx context.Context) ([]*entity.Order, error) {
	return s.orderRepo.GetOrdersForProcessing(ctx)
}

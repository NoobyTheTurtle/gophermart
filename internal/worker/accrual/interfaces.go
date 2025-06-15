package accrual

import (
	"context"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
)

type Logger interface {
	Info(message string, args ...any)
	Error(message string, args ...any)
}

type AccrualUseCase interface {
	ProcessOrder(ctx context.Context, orderNumber string) error
	GetOrdersForProcessing(ctx context.Context) ([]*entity.Order, error)
}

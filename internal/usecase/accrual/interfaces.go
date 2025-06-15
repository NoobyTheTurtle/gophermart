package accrual

import (
	"context"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
	"github.com/NoobyTheTurtle/gophermart/internal/repo/api/accrual"
	"github.com/NoobyTheTurtle/gophermart/internal/repo/postgres/balance"
	"github.com/NoobyTheTurtle/gophermart/internal/repo/postgres/order"
)

type AccrualAPI interface {
	GetOrderAccrual(ctx context.Context, orderNumber string) (*entity.Accrual, error)
}

var _ AccrualAPI = (*accrual.AccrualAPI)(nil)

type OrderRepo interface {
	GetOrdersForProcessing(ctx context.Context) ([]*entity.Order, error)
	GetOrderByNumber(ctx context.Context, number string) (*entity.Order, error)
	UpdateOrderStatus(ctx context.Context, number string, status entity.OrderStatus, accrual float64) error
}

var _ OrderRepo = (*order.OrderPostgresRepo)(nil)

type BalanceRepo interface {
	AddAccrual(ctx context.Context, userID int, accrual float64) error
}

var _ BalanceRepo = (*balance.BalancePostgresRepo)(nil)

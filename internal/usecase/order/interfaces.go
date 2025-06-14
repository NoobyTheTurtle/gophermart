package order

import (
	"context"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
	"github.com/NoobyTheTurtle/gophermart/internal/repo/postgres/balance"
	"github.com/NoobyTheTurtle/gophermart/internal/repo/postgres/order"
	"github.com/NoobyTheTurtle/gophermart/pkg/luhn"
)

type OrderRepo interface {
	CreateOrder(ctx context.Context, order *entity.Order) error
	GetOrderByNumber(ctx context.Context, number string) (*entity.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int) ([]*entity.Order, error)
	UpdateOrderStatus(ctx context.Context, number string, status entity.OrderStatus, accrual float64) error
	GetOrdersForProcessing(ctx context.Context) ([]*entity.Order, error)
}

var _ OrderRepo = (*order.OrderPostgresRepo)(nil)

type BalanceRepo interface {
	CreateBalance(ctx context.Context, userID int) error
	AddAccrual(ctx context.Context, userID int, accrual float64) error
}

var _ BalanceRepo = (*balance.BalancePostgresRepo)(nil)

type LuhnService interface {
	IsValid(number string) bool
}

var _ LuhnService = (*luhn.LuhnService)(nil)

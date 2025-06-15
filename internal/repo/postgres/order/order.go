package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type OrderPostgresRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *OrderPostgresRepo {
	return &OrderPostgresRepo{db: db}
}

func (r *OrderPostgresRepo) CreateOrder(ctx context.Context, order *entity.Order) error {
	err := r.db.QueryRowxContext(ctx, createOrderQuery, order.UserID, order.Number, order.Status, order.UploadedAt).Scan(&order.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entity.ErrOrderAlreadyExists
		}
		return fmt.Errorf("repo - postgres - order - CreateOrder: failed to create order: %w", err)
	}

	return nil
}

func (r *OrderPostgresRepo) GetOrderByNumber(ctx context.Context, number string) (*entity.Order, error) {
	order := &entity.Order{}

	err := r.db.GetContext(ctx, order, getOrderByNumberQuery, number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrOrderNotFound
		}
		return nil, fmt.Errorf("repo - postgres - order - GetOrderByNumber: failed to get order by number: %w", err)
	}

	return order, nil
}

func (r *OrderPostgresRepo) GetOrdersByUserID(ctx context.Context, userID int) ([]*entity.Order, error) {
	var orders []*entity.Order

	err := r.db.SelectContext(ctx, &orders, getOrdersByUserIDQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("repo - postgres - order - GetOrdersByUserID: failed to get orders by user ID: %w", err)
	}

	return orders, nil
}

func (r *OrderPostgresRepo) UpdateOrderStatus(ctx context.Context, number string, status entity.OrderStatus, accrual float64) error {
	_, err := r.db.ExecContext(ctx, updateOrderStatusQuery, status, accrual, number)
	if err != nil {
		return fmt.Errorf("repo - postgres - order - UpdateOrderStatus: failed to update order status: %w", err)
	}

	return nil
}

func (r *OrderPostgresRepo) GetOrdersForProcessing(ctx context.Context) ([]*entity.Order, error) {
	var orders []*entity.Order

	err := r.db.SelectContext(ctx, &orders, getOrdersForProcessingQuery)
	if err != nil {
		return nil, fmt.Errorf("repo - postgres - order - GetOrdersForProcessing: failed to get orders for processing: %w", err)
	}

	return orders, nil
}

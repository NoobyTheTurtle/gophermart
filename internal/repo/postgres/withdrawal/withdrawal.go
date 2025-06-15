package withdrawal

import (
	"context"
	"fmt"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
	"github.com/jmoiron/sqlx"
)

type WithdrawalPostgresRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *WithdrawalPostgresRepo {
	return &WithdrawalPostgresRepo{
		db: db,
	}
}

func (r *WithdrawalPostgresRepo) CreateWithdrawal(ctx context.Context, withdrawal *entity.Withdrawal) error {
	err := r.db.QueryRowxContext(ctx, createWithdrawalQuery,
		withdrawal.UserID,
		withdrawal.OrderNumber,
		withdrawal.Sum,
		withdrawal.ProcessedAt,
	).Scan(&withdrawal.ID)
	if err != nil {
		return fmt.Errorf("repo - withdrawal - CreateWithdrawal: %w", err)
	}

	return nil
}

func (r *WithdrawalPostgresRepo) GetWithdrawalsByUserID(ctx context.Context, userID int) ([]*entity.Withdrawal, error) {
	var withdrawals []*entity.Withdrawal

	err := r.db.SelectContext(ctx, &withdrawals, getWithdrawalsByUserIDQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("repo - withdrawal - GetWithdrawalsByUserID: %w", err)
	}

	return withdrawals, nil
}

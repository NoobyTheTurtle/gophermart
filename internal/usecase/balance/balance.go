package balance

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
)

type BalanceUseCaseImpl struct {
	balanceRepo BalanceRepo
}

func New(balanceRepo BalanceRepo) *BalanceUseCaseImpl {
	return &BalanceUseCaseImpl{
		balanceRepo: balanceRepo,
	}
}

func (uc *BalanceUseCaseImpl) GetUserBalance(ctx context.Context, userID int) (*entity.UserBalance, error) {
	userBalance, err := uc.balanceRepo.GetBalanceByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, entity.ErrBalanceNotFound) {
			if createErr := uc.balanceRepo.CreateBalance(ctx, userID); createErr != nil {
				return nil, fmt.Errorf("usecase - balance - GetUserBalance: failed to create balance: %w", createErr)
			}
			return &entity.UserBalance{
				UserID:    userID,
				Current:   0,
				Withdrawn: 0,
				UpdatedAt: time.Now(),
			}, nil
		}
		return nil, fmt.Errorf("usecase - balance - GetUserBalance: failed to get user balance: %w", err)
	}

	return userBalance, nil
}

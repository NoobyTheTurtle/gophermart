package withdrawal

import (
	"context"
	"fmt"
	"time"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
)

type WithdrawalUseCaseImpl struct {
	withdrawalRepo WithdrawalRepo
	balanceRepo    BalanceRepo
	luhnService    LuhnService
}

func New(withdrawalRepo WithdrawalRepo, balanceRepo BalanceRepo, luhnService LuhnService) *WithdrawalUseCaseImpl {
	return &WithdrawalUseCaseImpl{
		withdrawalRepo: withdrawalRepo,
		balanceRepo:    balanceRepo,
		luhnService:    luhnService,
	}
}

func (uc *WithdrawalUseCaseImpl) ProcessWithdrawal(ctx context.Context, userID int, orderNumber string, sum float64) error {
	if !uc.luhnService.IsValid(orderNumber) {
		return entity.ErrInvalidOrderNumber
	}

	balance, err := uc.balanceRepo.GetBalanceByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("usecase - withdrawal - ProcessWithdrawal: failed to get balance: %w", err)
	}

	if balance.Current < sum {
		return entity.ErrInsufficientFunds
	}

	withdrawal := &entity.Withdrawal{
		UserID:      userID,
		OrderNumber: orderNumber,
		Sum:         sum,
		ProcessedAt: time.Now(),
	}

	if err := uc.withdrawalRepo.CreateWithdrawal(ctx, withdrawal); err != nil {
		return fmt.Errorf("usecase - withdrawal - ProcessWithdrawal: failed to create withdrawal: %w", err)
	}

	balance.Current -= sum
	balance.Withdrawn += sum

	if err := uc.balanceRepo.UpdateBalance(ctx, balance); err != nil {
		return fmt.Errorf("usecase - withdrawal - ProcessWithdrawal: failed to update balance: %w", err)
	}

	return nil
}

func (uc *WithdrawalUseCaseImpl) GetUserWithdrawals(ctx context.Context, userID int) ([]*entity.Withdrawal, error) {
	withdrawals, err := uc.withdrawalRepo.GetWithdrawalsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("usecase - withdrawal - GetUserWithdrawals: failed to get withdrawals: %w", err)
	}

	return withdrawals, nil
}

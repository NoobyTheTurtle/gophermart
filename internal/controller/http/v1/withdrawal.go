package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/middleware"
	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/v1/request"
	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/v1/response"
	"github.com/NoobyTheTurtle/gophermart/internal/entity"
)

type withdrawalRoutes struct {
	withdrawalUseCase WithdrawalUseCase
	logger            Logger
}

func newWithdrawalRoutes(g *gin.RouterGroup, withdrawalUseCase WithdrawalUseCase, logger Logger) {
	r := &withdrawalRoutes{
		withdrawalUseCase: withdrawalUseCase,
		logger:            logger,
	}

	g.POST("/balance/withdraw", r.processWithdrawal)
	g.GET("/withdrawals", r.getUserWithdrawals)
}

// processWithdrawal handles withdrawal request
// @Summary Process withdrawal
// @Description Request withdrawal of loyalty points for order payment
// @Tags Withdrawal
// @Accept json
// @Produce json
// @Param request body request.Withdrawal true "Withdrawal request"
// @Success 200 "Withdrawal processed successfully"
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Error
// @Failure 402 {object} response.Error
// @Failure 422 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /api/user/balance/withdraw [post]
// @Security Bearer
func (r *withdrawalRoutes) processWithdrawal(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	var req request.Withdrawal
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request format")
		return
	}

	err := r.withdrawalUseCase.ProcessWithdrawal(c.Request.Context(), userID, req.Order, req.Sum)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidOrderNumber):
			response.ErrorResponse(c, http.StatusUnprocessableEntity, "invalid order number format")
			return
		case errors.Is(err, entity.ErrInsufficientFunds):
			response.ErrorResponse(c, http.StatusPaymentRequired, "insufficient funds")
			return
		default:
			r.logger.Error("Failed to process withdrawal", "userID", userID, "error", err)
			response.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	c.Status(http.StatusOK)
}

// getUserWithdrawals handles getting user withdrawals
// @Summary Get user withdrawals
// @Description Get list of user's withdrawals sorted by processing time desc
// @Tags Withdrawal
// @Produce json
// @Success 200 {array} response.Withdrawal
// @Success 204 "No withdrawals found"
// @Failure 401 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /api/user/withdrawals [get]
// @Security Bearer
func (r *withdrawalRoutes) getUserWithdrawals(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	withdrawals, err := r.withdrawalUseCase.GetUserWithdrawals(c.Request.Context(), userID)
	if err != nil {
		r.logger.Error("Failed to get user withdrawals", "userID", userID, "error", err)
		response.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if len(withdrawals) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	responseWithdrawals := make([]response.Withdrawal, len(withdrawals))
	for i, withdrawal := range withdrawals {
		responseWithdrawals[i] = response.Withdrawal{
			Order:       withdrawal.OrderNumber,
			Sum:         withdrawal.Sum,
			ProcessedAt: withdrawal.ProcessedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, responseWithdrawals)
}

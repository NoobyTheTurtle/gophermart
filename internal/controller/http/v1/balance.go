package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/middleware"
	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/v1/response"
)

type balanceRoutes struct {
	balanceUseCase BalanceUseCase
}

func newBalanceRoutes(g *gin.RouterGroup, balanceUseCase BalanceUseCase) {
	r := &balanceRoutes{
		balanceUseCase: balanceUseCase,
	}

	g.GET("/balance", r.getUserBalance)
}

// getUserBalance handles getting user balance
// @Summary Get user balance
// @Description Get current balance and total withdrawn amount
// @Tags Balance
// @Produce json
// @Success 200 {object} response.Balance
// @Failure 401 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /api/user/balance [get]
// @Security Bearer
func (r *balanceRoutes) getUserBalance(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	balance, err := r.balanceUseCase.GetUserBalance(c.Request.Context(), userID)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, response.Balance{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	})
}

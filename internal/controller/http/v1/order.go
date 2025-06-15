package v1

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/middleware"
	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/v1/response"
	"github.com/NoobyTheTurtle/gophermart/internal/entity"
)

type orderRoutes struct {
	orderUseCase OrderUseCase
	logger       Logger
}

func newOrderRoutes(g *gin.RouterGroup, orderUseCase OrderUseCase, logger Logger) {
	r := &orderRoutes{
		orderUseCase: orderUseCase,
		logger:       logger,
	}

	g.POST("/orders", r.uploadOrder)
	g.GET("/orders", r.getUserOrders)
}

// uploadOrder handles order upload
// @Summary Upload order number
// @Description Upload order number for loyalty points calculation
// @Tags Orders
// @Accept plain
// @Produce json
// @Param order body string true "Order number"
// @Success 200 "Order already uploaded by this user"
// @Success 202 "New order accepted for processing"
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Error
// @Failure 409 {object} response.Error
// @Failure 422 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /api/user/orders [post]
// @Security Bearer
func (r *orderRoutes) uploadOrder(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request format")
		return
	}

	orderNumber := strings.TrimSpace(string(body))
	if orderNumber == "" {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request format")
		return
	}

	err = r.orderUseCase.UploadOrder(c.Request.Context(), userID, orderNumber)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidOrderNumber):
			response.ErrorResponse(c, http.StatusUnprocessableEntity, "invalid order number format")
			return
		case errors.Is(err, entity.ErrOrderAlreadyExistsByAnotherUser):
			response.ErrorResponse(c, http.StatusConflict, "order already uploaded by another user")
			return
		case errors.Is(err, entity.ErrOrderAlreadyExists):
			c.Status(http.StatusOK)
			return
		default:
			r.logger.Error("Failed to upload order", "userID", userID, "orderNumber", orderNumber, "error", err)
			response.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	c.Status(http.StatusAccepted)
}

// getUserOrders handles getting user orders
// @Summary Get user orders
// @Description Get list of user's uploaded orders with statuses and accruals
// @Tags Orders
// @Produce json
// @Success 200 {array} response.Order
// @Success 204 "No orders found"
// @Failure 401 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /api/user/orders [get]
// @Security Bearer
func (r *orderRoutes) getUserOrders(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	orders, err := r.orderUseCase.GetUserOrders(c.Request.Context(), userID)
	if err != nil {
		r.logger.Error("Failed to get user orders", "userID", userID, "error", err)
		response.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if len(orders) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	responseOrders := make([]response.Order, len(orders))
	for i, order := range orders {
		responseOrders[i] = response.Order{
			Number:     order.Number,
			Status:     order.Status,
			UploadedAt: order.UploadedAt,
		}
		if order.Accrual > 0 {
			responseOrders[i].Accrual = &order.Accrual
		}
	}

	c.JSON(http.StatusOK, responseOrders)
}

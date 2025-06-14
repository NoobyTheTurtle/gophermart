package response

import (
	"time"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
)

type Order struct {
	Number     string             `json:"number" binding:"required" example:"1234567890"`
	Status     entity.OrderStatus `json:"status" binding:"required" example:"NEW"`
	Accrual    *float64           `json:"accrual,omitempty" binding:"required" example:"100.00"`
	UploadedAt time.Time          `json:"uploaded_at" binding:"required" example:"2021-01-01T00:00:00Z"`
}

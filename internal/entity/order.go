package entity

import (
	"errors"
	"time"
)

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID         int         `json:"id" db:"id"`
	UserID     int         `json:"user_id" db:"user_id"`
	Number     string      `json:"number" db:"number"`
	Status     OrderStatus `json:"status" db:"status"`
	Accrual    float64     `json:"accrual,omitempty" db:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at" db:"uploaded_at"`
}

var (
	ErrInvalidOrderNumber              = errors.New("invalid order number")
	ErrOrderAlreadyExists              = errors.New("order already exists")
	ErrOrderAlreadyExistsByAnotherUser = errors.New("order already exists by another user")
	ErrOrderNotFound                   = errors.New("order not found")
)

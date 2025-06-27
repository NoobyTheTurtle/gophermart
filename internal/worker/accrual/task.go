package accrual

import "github.com/NoobyTheTurtle/gophermart/internal/entity"

type orderTask struct {
	Order   *entity.Order
	Attempt int
}

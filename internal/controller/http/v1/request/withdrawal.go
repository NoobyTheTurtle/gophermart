package request

type Withdrawal struct {
	Order string  `json:"order" binding:"required" example:"2377225624"`
	Sum   float64 `json:"sum" binding:"required" example:"751"`
}

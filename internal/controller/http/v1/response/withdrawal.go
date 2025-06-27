package response

type Withdrawal struct {
	Order       string  `json:"order" binding:"required" example:"2377225624"`
	Sum         float64 `json:"sum" binding:"required" example:"500"`
	ProcessedAt string  `json:"processed_at" binding:"required" example:"2020-12-09T16:09:57+03:00"`
}

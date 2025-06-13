package response

type Ping struct {
	Status string `json:"status" binding:"required" example:"ok"`
}

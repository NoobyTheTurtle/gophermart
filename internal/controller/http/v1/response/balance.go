package response

type Balance struct {
	Current   float64 `json:"current" binding:"required" example:"100.00"`
	Withdrawn float64 `json:"withdrawn" binding:"required" example:"10.00"`
}

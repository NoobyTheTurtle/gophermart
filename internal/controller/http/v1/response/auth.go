package response

type Auth struct {
	Token string `json:"token" binding:"required" example:"token"`
}

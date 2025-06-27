package request

type Auth struct {
	Login    string `json:"login" binding:"required" example:"user"`
	Password string `json:"password" binding:"required" example:"password"`
}

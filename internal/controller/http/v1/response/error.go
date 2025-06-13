package response

type Error struct {
	Error string `json:"error" binding:"required" example:"message"`
}

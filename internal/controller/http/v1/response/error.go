package response

import "github.com/gin-gonic/gin"

type Error struct {
	Error string `json:"error" binding:"required" example:"message"`
}

func ErrorResponse(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, Error{Error: msg})
}

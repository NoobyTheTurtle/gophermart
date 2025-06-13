package v1

import (
	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/v1/response"
	"github.com/gin-gonic/gin"
)

// ping handles health check endpoint
// @Summary Health check
// @Description Check if the service is alive
// @Tags Health
// @Produce json
// @Success 200 {object} response.Ping
// @Router /ping [get]
func ping(c *gin.Context) {
	c.JSON(200, response.Ping{Status: "ok"})
}

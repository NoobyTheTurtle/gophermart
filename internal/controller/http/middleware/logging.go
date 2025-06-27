package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(log Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()
		size := c.Writer.Size()

		log.Info("HTTP request processed",
			"method", method,
			"path", path,
			"status", statusCode,
			"duration", duration,
			"size", size,
		)

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Error("Request error",
					"method", method,
					"path", path,
					"error", err.Err,
				)
			}
		}
	}
}

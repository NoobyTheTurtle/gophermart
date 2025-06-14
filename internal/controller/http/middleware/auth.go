package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/v1/response"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
	userIDCtxKey        = "userID"
)

func AuthMiddleware(authUseCase AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			response.ErrorResponse(c, http.StatusUnauthorized, "authentication required")
			c.Abort()
			return
		}

		userID, err := authUseCase.ValidateToken(token)
		if err != nil {
			response.ErrorResponse(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		c.Set(userIDCtxKey, userID)
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader(authorizationHeader)
	if authHeader != "" && strings.HasPrefix(authHeader, bearerPrefix) {
		return strings.TrimPrefix(authHeader, bearerPrefix)
	}

	token, err := c.Cookie("token")
	if err == nil {
		return token
	}

	return ""
}

func GetUserID(c *gin.Context) (int, bool) {
	userID, exists := c.Get(userIDCtxKey)
	if !exists {
		return 0, false
	}

	id, ok := userID.(int)
	return id, ok
}

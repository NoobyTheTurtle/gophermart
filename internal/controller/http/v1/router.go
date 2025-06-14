package v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/middleware"
)

func New(authUseCase AuthUseCase, orderUseCase OrderUseCase, balanceUseCase BalanceUseCase) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.GzipMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", ping)

	apiV1 := router.Group("/api/user")
	{
		newAuthRoutes(apiV1, authUseCase)

		protected := apiV1.Group("")
		protected.Use(middleware.AuthMiddleware(authUseCase))
		{
			newOrderRoutes(protected, orderUseCase)
			newBalanceRoutes(protected, balanceUseCase)
		}
	}

	return router
}

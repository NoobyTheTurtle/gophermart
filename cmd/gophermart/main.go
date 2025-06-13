package main

import (
	"context"

	"github.com/NoobyTheTurtle/gophermart/internal/app"
)

// @title           Gophermart API
// @version         1.0
// @description     API для накопительной системы лояльности «Гофермарт»

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	ctx := context.Background()

	app.Run(ctx)
}

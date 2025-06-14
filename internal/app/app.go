package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NoobyTheTurtle/gophermart/config"
	v1 "github.com/NoobyTheTurtle/gophermart/internal/controller/http/v1"
	balanceRepo "github.com/NoobyTheTurtle/gophermart/internal/repo/postgres/balance"
	orderRepo "github.com/NoobyTheTurtle/gophermart/internal/repo/postgres/order"
	userRepo "github.com/NoobyTheTurtle/gophermart/internal/repo/postgres/user"
	authUseCase "github.com/NoobyTheTurtle/gophermart/internal/usecase/auth"
	balanceUseCase "github.com/NoobyTheTurtle/gophermart/internal/usecase/balance"
	orderUseCase "github.com/NoobyTheTurtle/gophermart/internal/usecase/order"
	"github.com/NoobyTheTurtle/gophermart/pkg/jwt"
	"github.com/NoobyTheTurtle/gophermart/pkg/luhn"
	"github.com/NoobyTheTurtle/gophermart/pkg/password"
	"github.com/NoobyTheTurtle/gophermart/pkg/postgres"
)

func Run(ctx context.Context) {
	cfg := config.New()

	db, err := postgres.New(ctx, cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer postgres.Close(db)

	if err := runMigrations(db.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize services
	tokenService := jwt.New(cfg.JWTSecret, time.Hour*24)
	passwordService := password.New(password.DefaultCost)
	luhnService := luhn.New()

	// Initialize repositories
	userRepo := userRepo.New(db)
	orderRepo := orderRepo.New(db)
	balanceRepo := balanceRepo.New(db)

	// Initialize use cases
	authUseCaseImpl := authUseCase.New(userRepo, tokenService, passwordService)
	balanceUseCaseImpl := balanceUseCase.New(balanceRepo)
	orderUseCaseImpl := orderUseCase.New(orderRepo, balanceRepo, luhnService)

	router := v1.New(authUseCaseImpl, orderUseCaseImpl, balanceUseCaseImpl)

	server := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on %s", cfg.RunAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

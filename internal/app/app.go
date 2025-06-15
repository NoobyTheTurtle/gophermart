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
	accrualAPI "github.com/NoobyTheTurtle/gophermart/internal/repo/api/accrual"
	balanceRepo "github.com/NoobyTheTurtle/gophermart/internal/repo/postgres/balance"
	orderRepo "github.com/NoobyTheTurtle/gophermart/internal/repo/postgres/order"
	userRepo "github.com/NoobyTheTurtle/gophermart/internal/repo/postgres/user"
	withdrawalRepo "github.com/NoobyTheTurtle/gophermart/internal/repo/postgres/withdrawal"
	accrualUseCase "github.com/NoobyTheTurtle/gophermart/internal/usecase/accrual"
	authUseCase "github.com/NoobyTheTurtle/gophermart/internal/usecase/auth"
	balanceUseCase "github.com/NoobyTheTurtle/gophermart/internal/usecase/balance"
	orderUseCase "github.com/NoobyTheTurtle/gophermart/internal/usecase/order"
	withdrawalUseCase "github.com/NoobyTheTurtle/gophermart/internal/usecase/withdrawal"
	accrualWorker "github.com/NoobyTheTurtle/gophermart/internal/worker/accrual"
	"github.com/NoobyTheTurtle/gophermart/pkg/jwt"
	"github.com/NoobyTheTurtle/gophermart/pkg/logger"
	"github.com/NoobyTheTurtle/gophermart/pkg/luhn"
	"github.com/NoobyTheTurtle/gophermart/pkg/password"
	"github.com/NoobyTheTurtle/gophermart/pkg/postgres"
)

func Run(ctx context.Context) {
	cfg := config.New()

	zapLogger, err := logger.New()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	db, err := postgres.New(ctx, cfg.DatabaseURI)
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", "error", err)
	}

	defer func() {
		if errClose := postgres.Close(db); errClose != nil {
			zapLogger.Warn("Failed to close database", "error", errClose)
		}
	}()

	if err := runMigrations(db.DB); err != nil {
		zapLogger.Fatal("Failed to run migrations", "error", err)
	}

	// Initialize services
	tokenService := jwt.New(cfg.JWTSecret, time.Hour*24)
	passwordService := password.New(password.DefaultCost)
	luhnService := luhn.New()

	// Initialize API
	accrualAPIImpl := accrualAPI.New(cfg.AccrualSystemAddress)

	// Initialize repositories
	userRepo := userRepo.New(db)
	orderRepo := orderRepo.New(db)
	balanceRepo := balanceRepo.New(db)
	withdrawalRepo := withdrawalRepo.New(db)

	// Initialize use cases
	authUseCaseImpl := authUseCase.New(userRepo, tokenService, passwordService)
	balanceUseCaseImpl := balanceUseCase.New(balanceRepo)
	orderUseCaseImpl := orderUseCase.New(orderRepo, balanceRepo, luhnService)
	withdrawalUseCaseImpl := withdrawalUseCase.New(withdrawalRepo, balanceRepo, luhnService)
	accrualUseCaseImpl := accrualUseCase.New(
		accrualAPIImpl,
		orderRepo,
		balanceRepo,
	)

	// Initialize workers
	accrualWorkerImpl := accrualWorker.New(accrualUseCaseImpl, zapLogger, cfg.WorkerCount, cfg.ProcessInterval)

	accrualWorkerImpl.Start(ctx)

	router := v1.New(authUseCaseImpl, orderUseCaseImpl, balanceUseCaseImpl, withdrawalUseCaseImpl, zapLogger)

	server := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: router,
	}

	go func() {
		zapLogger.Info("Starting HTTP server", "address", cfg.RunAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("Failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("Received shutdown signal, starting graceful shutdown...")

	accrualWorkerImpl.Stop()

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxWithTimeout); err != nil {
		zapLogger.Error("Server forced to shutdown",
			"error", err,
			"timeout", "5s")
	} else {
		zapLogger.Info("Server exited gracefully")
	}
}

package main

import (
	"auth/internal/config"
	"auth/internal/handler"
	"auth/internal/logging"
	"auth/internal/middleware"
	"auth/internal/models"
	"auth/internal/repository"
	"auth/internal/router"
	"auth/internal/server"
	"auth/internal/service"
	"auth/internal/storage"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/subosito/gotenv"
	"go.uber.org/zap"
)

func main() {
	if err := gotenv.Load(".env"); err != nil {
		fmt.Println(err)
		return
	}
	config.Init()
	cfg := config.NewConfig()
	dbConfig := config.NewDBConfig()
	redisConfig := config.NewRedisConfig()
	ctx := context.Background()
	prefix := "refresh"

	logging.InitLogger(cfg.LoggingMode)
	logger := logging.Logger

	logger.Info("Starting server... with", zap.String("host", cfg.Host), zap.String("port", cfg.Port))

	// DATABASE
	db, err := storage.NewDatabase(dbConfig, models.ModelsList)
	if err != nil {
		logger.Fatal("Error while creating database", zap.Error(err))
	}

	// REDIS
	redisClient := storage.NewRedisClient(redisConfig)

	accessTokenTTL := time.Duration(cfg.AccessTokenExpiration) * time.Second

	// USER REPOSITORY
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(redisClient, ctx, prefix, cfg.JWTSecretKey, accessTokenTTL, "auth-service")
	// SERVICES
	userService := service.NewUserService(userRepo)

	tokenService := service.NewTokenService(tokenRepo, ctx, accessTokenTTL)
	authService := service.NewAuthService(userService, tokenService)

	// HANDLER
	h := handler.NewHandler(userService, authService)

	// MIDDLEWARE
	authMiddleware := middleware.AuthMiddleware(tokenService)

	// ROUTER
	router := router.SetupRouter(h, authMiddleware)
	srv, err := server.NewServer(cfg, h, router)

	if err != nil {
		logger.Fatal("Error while starting server", zap.Error(err))
	}
	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Error while starting server", zap.Error(err))
		}
	}()

	logger.Info("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		logger.Fatal("Error while stopping server")
	}

	logger.Info("Server stopped")
}

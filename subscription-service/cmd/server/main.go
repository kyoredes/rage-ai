package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"subscription/internal/config"
	"subscription/internal/handler"
	"subscription/internal/logging"
	"subscription/internal/middleware"
	"subscription/internal/models"
	"subscription/internal/repository"
	"subscription/internal/router"
	"subscription/internal/server"
	"subscription/internal/service"
	"subscription/internal/storage"
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
	ctx := context.Background()

	logging.InitLogger(cfg.LoggingMode)
	logger := logging.Logger

	devConfig := config.NewDevConfig()
	dbConfig := config.NewDBConfig()

	logger.Info("Starting server... with", zap.String("host", cfg.Host), zap.String("port", cfg.Port))
	db, err := storage.NewDatabase(dbConfig, models.ModelsList)
	if err != nil {
		logger.Fatal("failed to create database", zap.Error(err))
	}

	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)

	h := handler.NewHandler(subscriptionService)
	serverAuthMiddleware := middleware.DevAuthMiddleware(devConfig)
	router := router.SetupRouter(h, serverAuthMiddleware)

	srv, err := server.NewServer(cfg, h, router)
	if err != nil {
		logger.Fatal("failed to create server", zap.Error(err))
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
		logger.Error("Error while stopping server", zap.Error(err))
	}

	logger.Info("Server stopped")
}

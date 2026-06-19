package main

import (
	"fmt"
	"subscription/cmd/migration/script"
	"subscription/internal/config"
	"subscription/internal/logging"
	"subscription/internal/models"
	"subscription/internal/storage"

	"github.com/subosito/gotenv"
	"go.uber.org/zap"
)

func main() {
	config.Init()
	cfg := config.NewConfig()

	logging.InitLogger(cfg.LoggingMode)
	logger := logging.Logger
	err := gotenv.Load(".env")

	if err != nil {
		logger.Fatal("Error loading .env file", zap.Error(err))
	}
	db, err := storage.NewDatabase(config.NewDBConfig(), models.ModelsList)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := script.RunMigrations(db, models.ModelsList); err != nil {
		fmt.Println(err)
		return
	}
}

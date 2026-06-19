package storage

import (
	"auth/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(DBConfig *config.DBConfig, models ...any) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(DBConfig.DBDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil

}

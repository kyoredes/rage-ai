package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Uuid      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string `gorm:"type:varchar(255);uniqueIndex"`
	Password  string `gorm:"type:varchar(255);not null"`

	Telegram *TelegramUser `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type TelegramUser struct {
	Uuid       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	TelegramID string `gorm:"type:string;uniqueIndex"`

	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	User   *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

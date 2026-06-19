package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	Uuid      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;uniqueIndex"`
	UserID    uuid.UUID `gorm:"not null;uniqueIndex"` // кто подписан
	CreatedAt time.Time `db:"created_at"`
	StartsAt  time.Time `db:"starts_at"`
	ExpiresAt time.Time `db:"expires_at"`
}

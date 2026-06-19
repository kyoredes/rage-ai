package dto

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionDTO struct {
	Uuid      uuid.UUID
	UserID    uuid.UUID
	StartsAt  time.Time
	ExpiresAt time.Time
}

type CreateSubscriptionDTO struct {
	UserID    uuid.UUID
	StartsAt  time.Time
	ExpiresAt time.Time
}

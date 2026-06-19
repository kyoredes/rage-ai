package dto

import "time"

type UserListItem struct {
	UserID     string
	Email      string
	TelegramID string
	CreatedAt  time.Time
}

type UserDetail struct {
	UserID     string
	Email      string
	TelegramID string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type AuthStats struct {
	TotalUsers  int64
	NewUsers7d  int64
}

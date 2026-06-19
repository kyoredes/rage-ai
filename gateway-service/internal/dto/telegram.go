package dto

type TelegramUserDTO struct {
	TelegramID string `json:"telegramID" binding:"required"`
}

type TelegramInfo struct {
	TelegramID string `json:"telegramID"`
	UserID     string `json:"userID"`
	DeviceID   string `json:"deviceID"`
}

type TelegramProfile struct {
	TelegramID string `json:"telegramID"`
	UserID     string `json:"userID"`
	Email      string `json:"email,omitempty"`
}

type TelegramSubscription struct {
	SubscriptionID string `json:"subscriptionID"`
	UserID         string `json:"userID"`
	StartsAt       int64  `json:"startsAt"`
	ExpiresAt      int64  `json:"expiresAt"`
}

type TelegramChatDTO struct {
	TelegramID string `json:"telegramID" binding:"required"`
	Prompt     string `json:"prompt" binding:"required"`
}

type TelegramChatResponse struct {
	TelegramID string `json:"telegramID"`
	Response   string `json:"response"`
}

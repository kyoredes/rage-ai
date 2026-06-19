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

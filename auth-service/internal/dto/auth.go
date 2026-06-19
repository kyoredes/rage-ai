package dto

type TelegramStartResult struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}

type TelegramProfileResult struct {
	UserID     string
	TelegramID string
	Email      string
}

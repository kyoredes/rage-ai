package dto

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginResponse struct {
	Token string `json:"token"`
}

type AdminUser struct {
	UserID     string `json:"userID"`
	Email      string `json:"email"`
	TelegramID string `json:"telegramID"`
	CreatedAt  int64  `json:"createdAt"`
}

type AdminUserList struct {
	Users []AdminUser `json:"users"`
	Total int32       `json:"total"`
}

type AdminUserDetail struct {
	UserID     string `json:"userID"`
	Email      string `json:"email"`
	TelegramID string `json:"telegramID"`
	CreatedAt  int64  `json:"createdAt"`
	UpdatedAt  int64  `json:"updatedAt"`
}

type UpdateUserRequest struct {
	Email string `json:"email"`
}

type AdminSubscription struct {
	SubscriptionID string `json:"subscriptionID"`
	UserID         string `json:"userID"`
	StartsAt       int64  `json:"startsAt"`
	ExpiresAt      int64  `json:"expiresAt"`
}

type AdminSubscriptionList struct {
	Subscriptions []AdminSubscription `json:"subscriptions"`
	Total         int32               `json:"total"`
}

type UpdateSubscriptionRequest struct {
	StartsAt  int64 `json:"startsAt" binding:"required"`
	ExpiresAt int64 `json:"expiresAt" binding:"required"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatHistory struct {
	TelegramID string        `json:"telegramID"`
	Messages   []ChatMessage `json:"messages"`
}

type ChatSession struct {
	TelegramID    string `json:"telegramID"`
	MessageCount  int32  `json:"messageCount"`
}

type ChatSessionList struct {
	Sessions []ChatSession `json:"sessions"`
	Total    int32         `json:"total"`
}

type LLMConfig struct {
	Model       string   `json:"model"`
	Temperature float64  `json:"temperature"`
	MaxTokens   int32    `json:"maxTokens"`
	Debug       bool     `json:"debug"`
	Provider    string   `json:"provider"`
	G4FModels   []string `json:"g4fModels"`
	UsesLiteLLM bool     `json:"usesLiteLLM"`
}

type SystemPrompt struct {
	Prompt        string `json:"prompt"`
	DefaultPrompt string `json:"defaultPrompt"`
	IsCustom      bool   `json:"isCustom"`
}

type UpdateSystemPromptRequest struct {
	Prompt string `json:"prompt"`
}

type AdminStats struct {
	Users struct {
		Total  int32 `json:"total"`
		New7d  int32 `json:"new7d"`
	} `json:"users"`
	Subscriptions struct {
		Total   int32 `json:"total"`
		Active  int32 `json:"active"`
		Expired int32 `json:"expired"`
	} `json:"subscriptions"`
	Chat struct {
		Sessions int32 `json:"sessions"`
	} `json:"chat"`
}

type ServiceStatus struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	LatencyMs int64  `json:"latencyMs"`
}

type ServicesStatusResponse struct {
	Services  []ServiceStatus `json:"services"`
	CheckedAt int64           `json:"checkedAt"`
}

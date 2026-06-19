package handler

import (
	"gateway/internal/service"
)

type Handler struct {
	Telegram *TelegramHandler
	Admin    *AdminHandler
}

func NewHandler(
	telegramService *service.TelegramService,
	adminService *service.AdminService,
) *Handler {
	return &Handler{
		Telegram: NewTelegramHandler(telegramService),
		Admin:    NewAdminHandler(adminService),
	}
}

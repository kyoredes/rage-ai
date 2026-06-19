package handler

import (
	"auth/internal/dto"
	"auth/internal/logging"
	"auth/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	logger := logging.Logger
	var request dto.UserDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Debug("Wrong request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong request",
		})
		return
	}
	tokens, err := h.authService.LoginUser(request.Email, request.Password, request.DeviceID)
	if err != nil {
		logger.Error("Error logging in user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error logging in user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"tokens": tokens,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	logger := logging.Logger
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug("Wrong request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong request",
		})
		return
	}

	tokens, err := h.authService.RefreshToken(req.RefreshToken, req.DeviceID)
	if err != nil {
		logger.Error("Error refreshing token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error refreshing token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"tokens": tokens,
	})
}

func (h *AuthHandler) StartTelegram(c *gin.Context) {
	logger := logging.Logger
	var request dto.TelegramUserDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Debug("Wrong request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong request",
		})
		return
	}

	tokens, err := h.authService.StartTelegram(request.TelegramID)
	if err != nil {
		logger.Error("Error starting telegram", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error starting telegram",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"tokens": tokens,
	})
}

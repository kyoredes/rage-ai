package service

import (
	"auth/exceptions"
	"auth/internal/dto"
	"auth/internal/logging"

	"go.uber.org/zap"
)

type AuthService struct {
	userService  *UserService
	tokenService *TokenService
}

func NewAuthService(
	userService *UserService,
	tokenService *TokenService,
) *AuthService {
	return &AuthService{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (s *AuthService) StartTelegramWithUser(TelegramID string) (*dto.TelegramStartResult, error) {
	logger := logging.Logger
	telegramUser, err := s.userService.GetOrCreateTelegramUser(TelegramID)
	if err != nil {
		logger.Error("Error getting or creating telegram user", zap.Error(err))
		return nil, err
	}
	accessToken, err := s.tokenService.CreateAccessToken(telegramUser.UserID, TelegramID)
	if err != nil {
		logger.Error("Error creating access token", zap.Error(err))
		return nil, exceptions.ErrAccessTokenNotCreated
	}

	refreshToken, err := s.tokenService.CreateRefreshToken(telegramUser.UserID)
	if err != nil {
		logger.Error("Error creating refresh token", zap.Error(err))
		return nil, exceptions.ErrRefreshTokenNotCreated
	}

	if err := s.tokenService.SaveRefreshToken(refreshToken, telegramUser.UserID, TelegramID); err != nil {
		logger.Error("Error saving refresh token", zap.Error(err))
		return nil, err
	}

	return &dto.TelegramStartResult{
		UserID:       telegramUser.UserID.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

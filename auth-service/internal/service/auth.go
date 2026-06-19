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

func (s *AuthService) LoginUser(email, password, deviceID string) (*dto.Tokens, error) {
	logger := logging.Logger

	user, err := s.userService.CheckPassword(email, password)
	if err != nil {
		logger.Error("Error checking password", zap.Error(err))
		return nil, exceptions.ErrInvalidCredentials
	}

	accessToken, err := s.tokenService.CreateAccessToken(user.Uuid, deviceID)
	if err != nil {
		logger.Error("Error creating access token", zap.Error(err))
		return nil, exceptions.ErrAccessTokenNotCreated
	}

	refreshToken, err := s.tokenService.CreateRefreshToken(user.Uuid)
	if err != nil {
		logger.Error("Error creating refresh token", zap.Error(err))
		return nil, exceptions.ErrRefreshTokenNotCreated
	}

	if err := s.tokenService.SaveRefreshToken(refreshToken, user.Uuid, deviceID); err != nil {
		logger.Error("Error saving refresh token", zap.Error(err))
		return nil, err
	}

	return &dto.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken, deviceID string) (*dto.Tokens, error) {
	logger := logging.Logger
	userID, err := s.tokenService.GetRefreshToken(refreshToken, deviceID)
	if err != nil {
		logger.Error("Error getting refresh token", zap.Error(err))
		return nil, exceptions.ErrRefreshTokenNotFound
	}

	if err := s.tokenService.DeleteRefreshToken(refreshToken, deviceID); err != nil {
		logger.Error("Error deleting refresh token", zap.Error(err))
		return nil, exceptions.ErrDeleteRefreshToken
	}

	newAccessToken, err := s.tokenService.CreateAccessToken(userID, deviceID)
	if err != nil {
		logger.Error("Error creating access token", zap.Error(err))
		return nil, exceptions.ErrAccessTokenNotCreated
	}

	newRefreshToken, err := s.tokenService.CreateRefreshToken(userID)
	if err != nil {
		logger.Error("Error creating refresh token", zap.Error(err))
		return nil, exceptions.ErrRefreshTokenNotCreated
	}

	if err := s.tokenService.SaveRefreshToken(newRefreshToken, userID, deviceID); err != nil {
		logger.Error("Error saving refresh token", zap.Error(err))
		return nil, err
	}
	return &dto.Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) StartTelegram(TelegramID string) (*dto.Tokens, error) {
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

	return &dto.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

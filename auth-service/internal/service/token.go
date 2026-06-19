package service

import (
	"auth/internal/repository"
	"time"

	"github.com/google/uuid"
)

type TokenService struct {
	tokenRepository *repository.TokenRepository
	refreshTTL      time.Duration
}

func NewTokenService(tokenRepository *repository.TokenRepository, refreshTTL time.Duration) *TokenService {
	return &TokenService{
		tokenRepository: tokenRepository,
		refreshTTL:      refreshTTL,
	}
}

func (s *TokenService) SaveRefreshToken(token string, userID uuid.UUID, deviceID string) error {
	return s.tokenRepository.SetRefreshToken(token, userID, deviceID, s.refreshTTL)
}

func (s *TokenService) CreateAccessToken(userID uuid.UUID, deviceId string) (string, error) {
	return s.tokenRepository.CreateAccessToken(userID, deviceId)
}

func (s *TokenService) CreateRefreshToken(userID uuid.UUID) (string, error) {
	return s.tokenRepository.CreateRefreshToken(userID)
}

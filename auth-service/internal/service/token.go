package service

import (
	"auth/internal/logging"
	"auth/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TokenService struct {
	tokenRepository *repository.TokenRepository
	ctx             context.Context
	ttl             time.Duration
}

func NewTokenService(tokenRepository *repository.TokenRepository, ctx context.Context, ttl time.Duration) *TokenService {
	return &TokenService{
		tokenRepository: tokenRepository,
		ctx:             ctx,
		ttl:             ttl,
	}
}

func (s *TokenService) SaveRefreshToken(token string, userID uuid.UUID, deviceID string) error {
	err := s.tokenRepository.SetRefreshToken(token, userID, deviceID, s.ttl)
	if err != nil {
		return err
	}
	return nil
}

func (s *TokenService) GetRefreshToken(token, deviceID string) (uuid.UUID, error) {
	return s.tokenRepository.GetRefreshToken(token, deviceID)
}

func (s *TokenService) DeleteRefreshToken(token, deviceID string) error {
	return s.tokenRepository.DeleteRefreshToken(token, deviceID)
}

func (s *TokenService) CreateAccessToken(userID uuid.UUID, deviceId string) (string, error) {
	return s.tokenRepository.CreateAccessToken(userID, deviceId)
}

func (s *TokenService) CreateRefreshToken(userID uuid.UUID) (string, error) {
	return s.tokenRepository.CreateRefreshToken(userID)
}

func (s *TokenService) ParseAccessToken(token string, deviceId string) (uuid.UUID, error) {
	logger := logging.Logger

	parsedToken, err := s.tokenRepository.ParseAccessToken(token)
	if err != nil {
		logger.Error("error parsing access token", zap.Error(err))
		return uuid.Nil, err
	}

	claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	if !ok || !parsedToken.Valid {
		logger.Error("invalid access token", zap.Error(err))
		return uuid.Nil, errors.New("invalid access token")
	}
	audience, err := claims.GetAudience()
	if err != nil {
		logger.Error("error getting audience", zap.Error(err))
		return uuid.Nil, err
	}
	if audience[0] != deviceId {
		logger.Error("invalid deviceId ", zap.Error(err))
		return uuid.Nil, errors.New("invalid device-Id")
	}
	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		logger.Error("error parsing uuid", zap.Error(err))
		return uuid.Nil, err
	}

	return userId, nil
}

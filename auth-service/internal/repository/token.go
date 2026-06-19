package repository

import (
	"auth/internal/logging"
	"auth/internal/storage"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TokenRepository struct {
	client       *storage.RedisClient
	ctx          context.Context
	prefix       string
	accessSecret string
	ttl          time.Duration
	issuer       string
}

func NewTokenRepository(client *storage.RedisClient, ctx context.Context, prefix string, accessSecret string, ttl time.Duration, issuer string) *TokenRepository {
	if accessSecret == "" {
		panic("accessSecret is empty")
	}
	if len(accessSecret) < 32 {
		panic("accessSecret is too short")
	}
	return &TokenRepository{
		client:       client,
		ctx:          ctx,
		prefix:       prefix,
		accessSecret: accessSecret,
		ttl:          ttl,
		issuer:       issuer,
	}
}

func (r *TokenRepository) SetRefreshToken(token string, userID uuid.UUID, deviceID string, ttl time.Duration) error {
	return r.client.Client.SetEx(r.ctx, r.prefix+":"+token+":"+deviceID, userID.String(), ttl).Err()
}

func (r *TokenRepository) GetRefreshToken(token, deviceID string) (uuid.UUID, error) {
	logger := logging.Logger
	value := r.prefix + ":" + token + ":" + deviceID
	res, err := r.client.Client.Get(r.ctx, value).Result()
	if err != nil {
		logger.Error("error getting refresh token from redis", zap.String("token", value), zap.Error(err))
		return uuid.Nil, err
	}
	userID, err := uuid.Parse(res)
	if err != nil {
		logger.Error("error parsing uuid", zap.String("token", value), zap.Error(err))
		return uuid.Nil, err
	}
	return userID, nil
}

func (r *TokenRepository) DeleteRefreshToken(token string, deviceID string) error {
	return r.client.Client.Del(r.ctx, r.prefix+":"+token+":"+deviceID).Err()
}

func (s *TokenRepository) CreateAccessToken(userID uuid.UUID, deviceId string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		Issuer:    s.issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Audience:  jwt.ClaimStrings{deviceId},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.accessSecret))
}

func (s *TokenRepository) CreateRefreshToken(userID uuid.UUID) (string, error) {
	return uuid.New().String(), nil
}

func (s *TokenRepository) ParseAccessToken(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, s.TokenKeyFunc)
}

func (s *TokenRepository) TokenKeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, jwt.ErrSignatureInvalid
	}
	return []byte(s.accessSecret), nil
}

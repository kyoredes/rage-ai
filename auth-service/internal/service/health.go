package service

import (
	"auth/internal/repository"
	"auth/internal/storage"
	"context"
)

type HealthService struct {
	userRepo *repository.UserRepository
	redis    *storage.RedisClient
}

func NewHealthService(userRepo *repository.UserRepository, redis *storage.RedisClient) *HealthService {
	return &HealthService{userRepo: userRepo, redis: redis}
}

func (s *HealthService) Check(ctx context.Context) (dbOk, redisOk bool) {
	if _, err := s.userRepo.CountUsers(); err == nil {
		dbOk = true
	}
	if err := s.redis.Client.Ping(ctx).Err(); err == nil {
		redisOk = true
	}
	return dbOk, redisOk
}

package service

import (
	"context"
	"subscription/internal/repository"
)

type HealthService struct {
	repo *repository.SubscriptionRepository
}

func NewHealthService(repo *repository.SubscriptionRepository) *HealthService {
	return &HealthService{repo: repo}
}

func (s *HealthService) Check(ctx context.Context) bool {
	_, err := s.repo.CountAll()
	return err == nil
}

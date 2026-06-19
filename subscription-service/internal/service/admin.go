package service

import (
	"subscription/internal/dto"
	"subscription/internal/exception"
	"subscription/internal/logging"
	"subscription/internal/repository"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AdminService struct {
	repo *repository.SubscriptionRepository
}

func NewAdminService(repo *repository.SubscriptionRepository) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) ListSubscriptions(page, limit int, status string) ([]dto.SubscriptionDTO, int64, error) {
	subs, total, err := s.repo.ListSubscriptions(page, limit, status)
	if err != nil {
		logging.Logger.Error("ListSubscriptions failed", zap.Error(err))
		return nil, 0, err
	}

	result := make([]dto.SubscriptionDTO, len(subs))
	for i, sub := range subs {
		result[i] = dto.SubscriptionDTO{
			Uuid:      sub.Uuid,
			UserID:    sub.UserID,
			StartsAt:  sub.StartsAt,
			ExpiresAt: sub.ExpiresAt,
		}
	}
	return result, total, nil
}

func (s *AdminService) UpdateSubscription(subID string, startsAt, expiresAt time.Time) (*dto.SubscriptionDTO, error) {
	id, err := uuid.Parse(subID)
	if err != nil {
		return nil, exception.ErrSubscriptionNotFound
	}

	if !expiresAt.After(startsAt) {
		return nil, exception.ErrInvalidSubscriptionDates
	}

	sub, err := s.repo.GetSubByUuid(id)
	if err != nil {
		return nil, err
	}

	sub.StartsAt = startsAt
	sub.ExpiresAt = expiresAt

	updated, err := s.repo.UpdateSub(sub)
	if err != nil {
		return nil, err
	}

	return &dto.SubscriptionDTO{
		Uuid:      updated.Uuid,
		UserID:    updated.UserID,
		StartsAt:  updated.StartsAt,
		ExpiresAt: updated.ExpiresAt,
	}, nil
}

func (s *AdminService) DeleteSubscription(subID string) error {
	id, err := uuid.Parse(subID)
	if err != nil {
		return exception.ErrSubscriptionNotFound
	}
	return s.repo.DeleteSub(id)
}

func (s *AdminService) GetStats() (*dto.SubscriptionStats, error) {
	total, err := s.repo.CountAll()
	if err != nil {
		return nil, err
	}
	active, err := s.repo.CountActive()
	if err != nil {
		return nil, err
	}
	expired, err := s.repo.CountExpired()
	if err != nil {
		return nil, err
	}

	return &dto.SubscriptionStats{
		Total:   total,
		Active:  active,
		Expired: expired,
	}, nil
}

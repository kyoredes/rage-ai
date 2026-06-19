package repository

import (
	"errors"
	"subscription/internal/exception"
	"subscription/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		db: db,
	}
}

func (r *SubscriptionRepository) CreateSub(sub *models.Subscription) (*models.Subscription, error) {
	err := r.db.Create(sub).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, exception.ErrSubscriptionAlreadyExists
		}
		return nil, err
	}
	return sub, nil
}

func (r *SubscriptionRepository) GetSubByUuid(id uuid.UUID) (*models.Subscription, error) {
	var sub models.Subscription
	err := r.db.Where("uuid = ?", id).First(&sub).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrSubscriptionNotFound
	}
	if err != nil {
		return nil, err
	}
	return &sub, nil
}
func (r *SubscriptionRepository) UpdateSub(sub *models.Subscription) (*models.Subscription, error) {
	err := r.db.Save(sub).Error
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (r *SubscriptionRepository) DeleteSub(id uuid.UUID) error {
	var sub models.Subscription
	return r.db.Delete(&sub, id).Error
}

func (r *SubscriptionRepository) GetSubByUserId(userId uuid.UUID) (*models.Subscription, error) {
	var sub models.Subscription
	err := r.db.Where("user_id = ?", userId).First(&sub).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrSubscriptionNotFound
	}
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *SubscriptionRepository) ListSubscriptions(page, limit int, status string) ([]models.Subscription, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := r.db.Model(&models.Subscription{})
	now := time.Now()

	switch status {
	case "active":
		query = query.Where("expires_at > ?", now)
	case "expired":
		query = query.Where("expires_at <= ?", now)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var subs []models.Subscription
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&subs).Error; err != nil {
		return nil, 0, err
	}
	return subs, total, nil
}

func (r *SubscriptionRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&models.Subscription{}).Count(&count).Error
	return count, err
}

func (r *SubscriptionRepository) CountActive() (int64, error) {
	var count int64
	err := r.db.Model(&models.Subscription{}).Where("expires_at > ?", time.Now()).Count(&count).Error
	return count, err
}

func (r *SubscriptionRepository) CountExpired() (int64, error) {
	var count int64
	err := r.db.Model(&models.Subscription{}).Where("expires_at <= ?", time.Now()).Count(&count).Error
	return count, err
}

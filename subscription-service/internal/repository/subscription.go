package repository

import (
	"errors"
	"subscription/internal/exception"
	"subscription/internal/models"

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

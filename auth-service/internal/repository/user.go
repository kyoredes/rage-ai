package repository

import (
	"auth/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) DeleteUser(id uuid.UUID) error {
	var user models.User
	return r.db.Delete(&user, id).Error
}

func (r *UserRepository) CreateTelegrameUser(telegramUser *models.TelegramUser) error {
	return r.db.Create(telegramUser).Error
}

func (r *UserRepository) GetTelegramUserByID(id uuid.UUID) (*models.TelegramUser, error) {
	var user models.TelegramUser

	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepository) GetTelegramUserByTelegramID(TelegramID string) (*models.TelegramUser, error) {
	var user models.TelegramUser

	if err := r.db.Where("telegram_id = ?", TelegramID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepository) UpdateTelegramUser(user *models.TelegramUser) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) DeleteTelegramUser(id uuid.UUID) error {
	var user models.TelegramUser
	return r.db.Delete(&user, id).Error
}

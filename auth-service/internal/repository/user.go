package repository

import (
	"auth/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserListItem struct {
	User       models.User
	TelegramID string
}

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

func (r *UserRepository) ListUsers(page, limit int, search string) ([]UserListItem, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	var total int64
	countQuery := r.db.Model(&models.User{}).
		Joins("LEFT JOIN telegram_users ON telegram_users.user_id = users.uuid")
	if search != "" {
		pattern := "%" + search + "%"
		countQuery = countQuery.Where("users.email ILIKE ? OR telegram_users.telegram_id ILIKE ?", pattern, pattern)
	}
	if err := countQuery.Distinct("users.uuid").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	listQuery := r.db.Model(&models.User{}).
		Select("users.*, COALESCE(telegram_users.telegram_id, '') as telegram_id").
		Joins("LEFT JOIN telegram_users ON telegram_users.user_id = users.uuid")
	if search != "" {
		pattern := "%" + search + "%"
		listQuery = listQuery.Where("users.email ILIKE ? OR telegram_users.telegram_id ILIKE ?", pattern, pattern)
	}

	type row struct {
		models.User
		TelegramID string `gorm:"column:telegram_id"`
	}
	var rows []row
	if err := listQuery.Order("users.created_at DESC").Offset(offset).Limit(limit).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]UserListItem, len(rows))
	for i, row := range rows {
		items[i] = UserListItem{User: row.User, TelegramID: row.TelegramID}
	}
	return items, total, nil
}

func (r *UserRepository) CountUsers() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}

func (r *UserRepository) CountUsersSince(since time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("created_at >= ?", since).Count(&count).Error
	return count, err
}

func (r *UserRepository) GetTelegramUserByUserID(userID uuid.UUID) (*models.TelegramUser, error) {
	var user models.TelegramUser
	if err := r.db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

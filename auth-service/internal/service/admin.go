package service

import (
	"auth/exceptions"
	"auth/internal/dto"
	"auth/internal/logging"
	"auth/internal/repository"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AdminService struct {
	repo *repository.UserRepository
}

func NewAdminService(repo *repository.UserRepository) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) ListUsers(page, limit int, search string) ([]dto.UserListItem, int64, error) {
	items, total, err := s.repo.ListUsers(page, limit, search)
	if err != nil {
		logging.Logger.Error("ListUsers failed", zap.Error(err))
		return nil, 0, err
	}

	result := make([]dto.UserListItem, len(items))
	for i, item := range items {
		email := item.User.Email
		if strings.HasSuffix(email, "@telegram.org") {
			email = ""
		}
		result[i] = dto.UserListItem{
			UserID:     item.User.Uuid.String(),
			Email:      email,
			TelegramID: item.TelegramID,
			CreatedAt:  item.User.CreatedAt,
		}
	}
	return result, total, nil
}

func (s *AdminService) GetUser(userID string) (*dto.UserDetail, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, exceptions.ErrUserNotFound
	}

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrUserNotFound
		}
		return nil, err
	}

	email := user.Email
	telegramID := ""
	tgUser, err := s.repo.GetTelegramUserByUserID(id)
	if err == nil {
		telegramID = tgUser.TelegramID
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if strings.HasSuffix(email, "@telegram.org") {
		email = ""
	}

	return &dto.UserDetail{
		UserID:     user.Uuid.String(),
		Email:      email,
		TelegramID: telegramID,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}, nil
}

func (s *AdminService) UpdateUser(userID, email string) (*dto.UserDetail, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, exceptions.ErrUserNotFound
	}

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrUserNotFound
		}
		return nil, err
	}

	if email != "" {
		user.Email = email
	}
	if err := s.repo.UpdateUser(user); err != nil {
		return nil, err
	}

	return s.GetUser(userID)
}

func (s *AdminService) DeleteUser(userID string) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return exceptions.ErrUserNotFound
	}

	tgUser, err := s.repo.GetTelegramUserByUserID(id)
	if err == nil {
		if err := s.repo.DeleteTelegramUser(tgUser.Uuid); err != nil {
			return err
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return s.repo.DeleteUser(id)
}

func (s *AdminService) GetStats() (*dto.AuthStats, error) {
	total, err := s.repo.CountUsers()
	if err != nil {
		return nil, err
	}

	since := time.Now().AddDate(0, 0, -7)
	new7d, err := s.repo.CountUsersSince(since)
	if err != nil {
		return nil, err
	}

	return &dto.AuthStats{
		TotalUsers: total,
		NewUsers7d: new7d,
	}, nil
}

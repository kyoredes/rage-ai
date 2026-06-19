package service

import (
	"auth/exceptions"
	"auth/internal/dto"
	"auth/internal/logging"
	"auth/internal/models"
	"auth/internal/repository"
	"auth/internal/security"
	"errors"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) RegisterUser(email, password string) (*models.User, error) {
	logger := logging.Logger
	_, err := s.repo.GetUserByEmail(email)
	if err == nil {
		logger.Error("User already exists", zap.Error(err))
		return nil, exceptions.ErrUserAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		logger.Error("Error hashing password", zap.Error(err))
		return nil, exceptions.ErrCreatingUser
	}
	user := &models.User{
		Email:    email,
		Password: hashedPassword,
	}
	err = s.repo.CreateUser(user)
	if err != nil {
		logger.Error("Error creating user", zap.Error(err))
		return nil, exceptions.ErrCreatingUser
	}
	return user, nil
}

func (s *UserService) CreateTelegrameUser(telegramUser *models.TelegramUser) error {
	return s.repo.CreateTelegrameUser(telegramUser)
}

func (s *UserService) GetOrCreateTelegramUser(TelegramID string) (*models.TelegramUser, error) {
	logger := logging.Logger
	existingTelegramUser, err := s.repo.GetTelegramUserByTelegramID(TelegramID)
	if err == nil {
		return existingTelegramUser, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Error getting telegram user", zap.Error(err))
		return nil, err
	}

	user, err := s.RegisterUser(
		TelegramID+"@telegram.org",
		uuid.NewString(),
	)
	if err != nil {
		return nil, err
	}

	telegramUser := &models.TelegramUser{
		UserID:     user.Uuid,
		User:       user,
		TelegramID: TelegramID,
	}

	if err := s.CreateTelegrameUser(telegramUser); err != nil {
		logger.Error("Error creating telegram user", zap.Error(err))
		return nil, err
	}

	return telegramUser, nil
}

func (s *UserService) GetTelegramProfile(telegramID string) (*dto.TelegramProfileResult, error) {
	logger := logging.Logger
	telegramUser, err := s.repo.GetTelegramUserByTelegramID(telegramID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrUserNotFound
		}
		logger.Error("Error getting telegram user", zap.Error(err))
		return nil, err
	}

	user, err := s.repo.GetUserByID(telegramUser.UserID)
	if err != nil {
		logger.Error("Error getting user for telegram profile", zap.Error(err))
		return nil, err
	}

	email := user.Email
	if strings.HasSuffix(email, "@telegram.org") {
		email = ""
	}

	return &dto.TelegramProfileResult{
		UserID:     user.Uuid.String(),
		TelegramID: telegramUser.TelegramID,
		Email:      email,
	}, nil
}

package service

import (
	"auth/exceptions"
	"auth/internal/logging"
	"auth/internal/models"
	"auth/internal/repository"
	"auth/internal/security"

	"github.com/google/uuid"
	"go.uber.org/zap"
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

func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *UserService) CheckPassword(email, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if err := security.CheckPassword(password, user.Password); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) CreateTelegrameUser(telegramUser *models.TelegramUser) error {
	return s.repo.CreateTelegrameUser(telegramUser)
}

func (s *UserService) GetOrCreateTelegramUser(TelegramID string) (*models.TelegramUser, error) {
	logger := logging.Logger
	existingTelegramUser, err := s.repo.GetTelegramUserByTelegramID(TelegramID)
	if existingTelegramUser != nil {
		return existingTelegramUser, nil
	}

	user, err := s.RegisterUser(
		TelegramID+"@telegram.org",
		uuid.NewString(),
	)
	telegramUser := &models.TelegramUser{
		UserID:     user.Uuid,
		User:       user,
		TelegramID: TelegramID,
	}

	err = s.CreateTelegrameUser(telegramUser)
	if err != nil {
		logger.Error("Error creating telegram user", zap.Error(err))
		return nil, err
	}

	return telegramUser, nil
}

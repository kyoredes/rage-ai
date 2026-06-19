package handler

import (
	"auth/exceptions"
	"auth/internal/dto"
	"auth/internal/logging"
	"auth/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	logger := logging.Logger
	var request dto.CreateUserDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Debug("Wrong request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong request",
		})
		return
	}

	user, err := h.service.RegisterUser(request.Email, request.Password)
	if err != nil {
		if err == exceptions.ErrUserAlreadyExists {
			logger.Error("User already exists", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User already exists",
			})
			return
		}
		if err == exceptions.ErrCreatingUser {
			logger.Error("Error creating user", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error creating user",
			})
			return
		}
		logger.Error("Error creating user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creating user",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "created",
		"email":  user.Email,
	})

}

func (h *UserHandler) GetUser(c *gin.Context) {
	userIdRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not found",
		})
		return
	}
	if userIdRaw == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not found",
		})
		return
	}
	userId, ok := userIdRaw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not found",
		})
		return
	}

	user, err := h.service.GetUserByID(userId)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not found",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"id":    user.Uuid,
		"email": user.Email,
	})
}

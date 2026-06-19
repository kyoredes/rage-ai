package handler

import (
	"net/http"
	"subscription/internal/dto"
	"subscription/internal/exception"
	"subscription/internal/logging"
	"subscription/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SubHandler struct {
	service *service.SubscriptionService
}

func NewSubHandler(service *service.SubscriptionService) *SubHandler {
	return &SubHandler{
		service: service,
	}
}

func (h *SubHandler) CreateSubscription(c *gin.Context) {
	logger := logging.Logger
	var request dto.CreateSubscriptionDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Debug("Wrong request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong request",
		})
		return
	}

	sub, err := h.service.GetOrCreateSub(&request)

	if err != nil {
		if err == exception.ErrSubscriptionAlreadyExists {
			logger.Error("Subscription already exists", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Subscription already exists",
			})
			return
		}
		if err == exception.ErrCreatingSubscription {
			logger.Error("Error creating subscription", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error creating subscription",
			})
			return
		}
		logger.Error("Error creating subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creating subscription",
		})
		return
	}
	c.JSON(http.StatusCreated, dto.Response{
		Status: "created",
		Data:   sub,
	})

}

func (h *SubHandler) GetSubscription(c *gin.Context) {
	logger := logging.Logger
	subUuidRaw := c.Query("sub_id")
	userIdRaw := c.Query("user_id")
	if subUuidRaw == "" && userIdRaw == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "sub_id or user_id should be provided",
		})
		return
	}
	if subUuidRaw != "" && userIdRaw != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "only one query parameter should be provided",
		})
		return
	}
	if subUuidRaw != "" {
		subUuid, err := uuid.Parse(subUuidRaw)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "subscription uuid is not valid",
			})
			return
		}

		sub, err := h.service.GetSubByUuid(subUuid)

		if err != nil {
			logger.Error("Error getting subscription", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "subscription not found",
			})
			return
		}
		c.JSON(http.StatusOK, dto.Response{
			Status: "ok",
			Data:   sub,
		})
		return
	}
	if userIdRaw != "" {
		userId, err := uuid.Parse(userIdRaw)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user uuid is not valid",
			})
			return
		}
		sub, err := h.service.GetSubByUserId(userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "subscription not found",
			})
			return
		}
		c.JSON(http.StatusOK, dto.Response{
			Status: "ok",
			Data:   sub,
		})
		return
	}
}

package router

import (
	"gateway/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handler.Handler, authMiddleware gin.HandlerFunc) *gin.Engine {
	router := gin.Default()

	tg := router.Group("/telegram")
	tg.Use(authMiddleware)
	tg.POST("/start", h.Telegram.StartTelegram)
	tg.POST("/profile", h.Telegram.GetProfile)
	tg.POST("/subscription", h.Telegram.GetSubscription)

	return router
}

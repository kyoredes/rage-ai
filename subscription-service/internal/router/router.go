package router

import (
	"subscription/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handler.Handler, authMiddleware gin.HandlerFunc) *gin.Engine {
	router := gin.Default()

	sub := router.Group("/subscription")
	sub.Use(authMiddleware)
	sub.POST("/create", h.Subscription.CreateSubscription)
	sub.GET("/", h.Subscription.GetSubscription)

	return router
}

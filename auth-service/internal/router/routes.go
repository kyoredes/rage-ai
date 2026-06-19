package router

import (
	"auth/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handler.Handler, authMiddleware gin.HandlerFunc) *gin.Engine {
	router := gin.Default()

	user := router.Group("/user")
	user.POST("/create", h.User.CreateUser)
	user.Use(authMiddleware)
	user.GET("/me", h.User.GetUser)

	auth := router.Group("/auth")
	auth.POST("/login", h.Auth.LoginUser)
	auth.POST("/refresh", h.Auth.RefreshToken)

	telegram := router.Group("/telegram")
	telegram.POST("/start", h.Auth.StartTelegram)

	return router
}

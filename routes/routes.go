package routes

import (
	"github.com/gin-gonic/gin"
	"order_service/cmd/handler"
	"order_service/middleware"
)

func SetupRoutes(router *gin.Engine, orderHandler handler.OrderHandler, jwtSecret string) {
	router.Use(middleware.RequestLogger())
	authMiddleware := middleware.AuthMiddleware(jwtSecret)
	router.Use(authMiddleware)
}

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
	router.POST("/v1/checkout", orderHandler.CheckOutOrder)
	router.GET("/v1/order_history", orderHandler.GetOrderHistory)
}

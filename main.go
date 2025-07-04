package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"order_service/cmd/handler"
	"order_service/cmd/repository"
	"order_service/cmd/resource"
	"order_service/cmd/service"
	"order_service/cmd/usecase"
	"order_service/config"
	"order_service/infra/log"
	"order_service/kafka"
	"order_service/routes"
)

func main() {
	cfg := config.LoadConfig(
		config.WithConfigFolder([]string{"./files/config"}),
		config.WithConfigFile("config"),
		config.WithConfigType("yaml"),
	)
	log.SetupLogger()

	db := resource.InitDB(&cfg)
	redis := resource.InitRedis(&cfg)
	kafkaProducer := kafka.NewKafkaProducer([]string{"localhost:9093"}, "order.created")

	orderRepo := repository.NewOrderRepository(db, redis, cfg.ProductService.Host)
	orderService := service.NewOrderService(*orderRepo)
	orderUseCase := usecase.NewOrderUseCase(*orderService, *kafkaProducer)
	orderHandler := handler.NewHandler(*orderUseCase)

	router := gin.Default()
	routes.SetupRoutes(router, *orderHandler, cfg.Secrete.JWTSecret)

	_ = router.Run(":" + cfg.App.Port)
	fmt.Println("Server running on port:", cfg.App.Port)
}

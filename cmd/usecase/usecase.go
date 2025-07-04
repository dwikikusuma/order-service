package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"order_service/cmd/service"
	"order_service/infra/constant"
	"order_service/infra/log"
	"order_service/kafka"
	"order_service/models"
	"time"
)

type OrderUseCase struct {
	OrderService  service.OrderService
	KafkaProducer kafka.KafkaProducer
}

func NewOrderUseCase(orderService service.OrderService, kafkaProducer kafka.KafkaProducer) *OrderUseCase {
	return &OrderUseCase{
		OrderService:  orderService,
		KafkaProducer: kafkaProducer,
	}
}

func (uc *OrderUseCase) CheckOutOrder(ctx context.Context, param *models.CheckoutRequest) (int64, error) {
	if param.IdempotencyToken != "" {
		isExists, err := uc.OrderService.CheckIdempotency(ctx, param.IdempotencyToken)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"idempotencyToken": param.IdempotencyToken,
				"error":            err,
			}).Error("Error checking idempotency")
			return 0, err
		}

		if isExists {
			return 0, fmt.Errorf("order with idempotency token '%s' already processed", param.IdempotencyToken)
		}
	}

	// validate products
	if err := uc.validateProducts(param.Items); err != nil {
		log.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error validating products")
		return 0, err
	}

	totalQty, totalAmount := uc.calculateOrderSummary(param.Items)
	products, orderHistory := uc.constructOrderDetail(param.Items)

	orderDetail := &models.OrderDetail{
		Products:     products,
		OrderHistory: orderHistory,
	}

	order := &models.Order{
		UserID:          param.UserID,
		Amount:          totalAmount,
		TotalQty:        int(totalQty),
		Status:          constant.OrderStatusCreated,
		PaymentMethod:   param.PaymentMethod,
		ShippingAddress: param.ShippingAddress,
	}

	// Save order and order detail, and handle idempotency token and Kafka event
	orderID, err := uc.OrderService.SaveOrderAndOrderDetail(ctx, order, orderDetail, param.IdempotencyToken, uc.KafkaProducer)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func (uc *OrderUseCase) validateProducts(items []models.CheckoutItem) error {
	seen := map[int64]bool{}
	for _, item := range items {
		if seen[item.ProductID] {
			return errors.New("duplicate product in checkout")
		}
		seen[item.ProductID] = true

		if item.Quantity <= 0 {
			return errors.New("quantity must be greater than zero")
		}

		if item.Price <= 0 {
			return errors.New("price must be greater than zero")
		}
	}
	return nil
}

func (uc *OrderUseCase) calculateOrderSummary(items []models.CheckoutItem) (int64, float64) {
	var totalQty int64
	var totalAmount float64

	for _, item := range items {
		totalQty += item.Quantity
		totalAmount += item.Price
	}
	return totalQty, totalAmount
}

func (uc *OrderUseCase) constructOrderDetail(items []models.CheckoutItem) (string, string) {
	orderDetail, _ := json.Marshal(items)
	orderHistory := []map[string]interface{}{
		{"status": "created", "timestamp": time.Now()},
	}

	orderHistoryJson, _ := json.Marshal(orderHistory)
	return string(orderDetail), string(orderHistoryJson)
}

func (uc *OrderUseCase) GetOrderHistoryByUserId(ctx context.Context, param *models.OrderHistoryParam) ([]models.OrderHistoryResponse, error) {
	orderHistories, err := uc.OrderService.GetOrderHistoryByUserId(ctx, param)
	if err != nil {
		return nil, err
	}
	return orderHistories, nil
}

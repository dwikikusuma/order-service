package usecase

import (
	"context"
	"encoding/json"
	"errors"
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
	var orderID int64

	// check idempotency
	if param.IdempotencyToken != "" {
		isExists, err := uc.OrderService.CheckIdempotency(ctx, param.IdempotencyToken)
		if err != nil {
			return 0, err
		}

		if isExists {
			return 0, errors.New("order already processed")
		}
	}

	// validate products
	if err := uc.validateProducts(param.Items); err != nil {
		return 0, err
	}

	// calculate product amount
	totalQty, totalAmount := uc.calculateOrderSummary(param.Items)

	// construct order and order detail
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

	// save order and order detail
	orderID, insertErr := uc.OrderService.SaveOrderAndOrderDetail(ctx, order, orderDetail)
	if insertErr != nil {
		log.Logger.WithFields(logrus.Fields{
			"message": "error occurred on uc.OrderService.SaveOrderAndOrderDetail",
			"error":   insertErr,
		}).Info("failed to save order and order detail")
		return 0, insertErr
	}

	// save idempotency
	if param.IdempotencyToken != "" {
		err := uc.OrderService.SaveIdempotency(ctx, param.IdempotencyToken)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"idempotency": param.IdempotencyToken,
				"error":       err,
			}).Info("error occurred on uc.OrderService.SaveIdempotency")
			return 0, err
		}
	}

	// TO DO:
	// connect payment service
	err := uc.KafkaProducer.PublishOrderCreated(ctx, &models.OrderCreatedEvent{
		OrderID:         orderID,
		UserID:          param.UserID,
		TotalAmount:     totalAmount,
		PaymentMethod:   param.PaymentMethod,
		ShippingAddress: param.ShippingAddress,
	})

	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"err":     err.Error(),
			"message": "failed to create order event on uc.KafkaProducer.PublishOrderCreated",
		})
		return 0, err
	}
	// checkout order -> Done
	// order history
	// connect to payment service
	// connect to product service -> validate product and validity management

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

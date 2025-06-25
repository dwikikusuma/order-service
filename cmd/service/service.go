package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"order_service/cmd/repository"
	"order_service/infra/log"
	"order_service/models"
)

type OrderService struct {
	OrderRepository repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) *OrderService {
	return &OrderService{
		OrderRepository: orderRepo,
	}
}

// CheckIdempotency check idempotency
func (s *OrderService) CheckIdempotency(ctx context.Context, idempotencyKey string) (bool, error) {
	isExists, err := s.OrderRepository.CheckIdempotency(ctx, idempotencyKey)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"message": "error occurred on s.OrderRepository.CheckIdempotency",
			"error":   err,
		})
		return false, err
	}
	return isExists, nil
}

// SaveIdempotency save idempotency
func (s *OrderService) SaveIdempotency(ctx context.Context, idempotencyKey string) error {
	err := s.OrderRepository.SaveIdempotency(ctx, idempotencyKey)
	if err != nil {
		return err
	}
	return nil
}

// SaveOrderAndOrderDetail save order and order_detail
func (s *OrderService) SaveOrderAndOrderDetail(ctx context.Context, order *models.Order, orderDetail *models.OrderDetail) (int64, error) {
	var orderID int64
	err := s.OrderRepository.WithTransaction(ctx, func(tx *gorm.DB) error {
		err := s.OrderRepository.InsertOrderDetailTx(ctx, tx, orderDetail)
		if err != nil {
			return err
		}

		order.OrderDetailID = orderDetail.ID
		err = s.OrderRepository.InsertOrderTx(ctx, tx, order)
		if err != nil {
			return err
		}

		orderID = order.ID
		return nil
	})

	if err != nil {
		return 0, err
	}

	return orderID, nil
}

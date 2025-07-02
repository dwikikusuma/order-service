package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"order_service/infra/constant"
	"order_service/infra/log"
	"order_service/models"
	"time"
)

// WithTransaction manages a database transaction and ensures that the
// transaction is either committed or rolled back based on the outcome
// of the provided callback function. It also handles panics gracefully
// to ensure that resources are cleaned up properly.
//
// The method begins a new transaction, passes it to the provided
// callback function `fn`, and ensures that:
//  1. The transaction is rolled back if `fn` returns an error.
//  2. The transaction is rolled back if a panic occurs inside `fn`.
//  3. The transaction is committed if `fn` executes successfully without error.
//
// Parameters:
//   - `ctx`: A context to associate with the transaction. It allows for
//     cancellation, timeouts, and request-scoped values.
//   - `fn`: A callback function that takes a `*gorm.DB` (the transaction)
//     and returns an `error`. This function will execute within the scope
//     of the transaction.
//
// Returns:
//   - An error if the transaction failed or if `fn` returned an error.
//   - `nil` if the transaction was successful and committed.
func (r *OrderRepository) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	// Begin a new transaction and associate it with the provided context.
	tx := r.Database.Begin().WithContext(ctx)

	// Defer a function to recover from any panics that may occur within the transaction.
	defer func() {
		// If a panic occurs, the recover function will handle it by rolling back the transaction
		// and re-raising the panic to allow higher-level handlers to catch it.
		if r := recover(); r != nil {
			tx.Rollback() // Rollback the transaction to ensure no partial data is committed.
			panic(r)      // Re-raise the panic so it can be handled elsewhere in the program.
		}
	}()

	// Execute the user-defined function `fn` with the transaction `tx` as an argument.
	// If `fn` returns an error, roll back the transaction and return the error.
	if err := fn(tx); err != nil {
		tx.Rollback() // Rollback the transaction if an error occurs in `fn`.
		return err    // Return the error to indicate failure.
	}

	// If `fn` executes successfully, commit the transaction and return any potential errors.
	return tx.Commit().Error // Return the result of the commit (or any error that occurred).
}

// InsertOrderTx insert order
func (r *OrderRepository) InsertOrderTx(ctx context.Context, tx *gorm.DB, order *models.Order) error {
	err := tx.WithContext(ctx).Table("orders").Create(&order).Error
	return err
}

// InsertOrderDetailTx insert order detail
func (r *OrderRepository) InsertOrderDetailTx(ctx context.Context, tx *gorm.DB, orderDetail *models.OrderDetail) error {
	err := tx.WithContext(ctx).Table("order_detail").Create(&orderDetail).Error
	return err
}

// CheckIdempotency check idempotency
func (r *OrderRepository) CheckIdempotency(ctx context.Context, idempotencyKey string) (bool, error) {
	var reqLog *models.OrderRequestLog
	err := r.Database.WithContext(ctx).Table("order_request_log").First(&reqLog, "idempotency_token = ?", idempotencyKey).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// SaveIdempotency save idempotency
func (r *OrderRepository) SaveIdempotency(ctx context.Context, idempotencyKey string) error {
	orderLog := models.OrderRequestLog{
		IdempotencyToken: idempotencyKey,
		CreateTime:       time.Now(),
	}

	err := r.Database.WithContext(ctx).Table("order_request_log").Create(&orderLog).Error
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"message": fmt.Sprintf("error occured on r.Database.WithContext(ctx).Table(\"order_request_log\").Create(&orderLog).Error"),
			"error":   err,
		})
		return err
	}
	return nil
}

func (r *OrderRepository) GetOrderHistoryByUserId(ctx context.Context, param *models.OrderHistoryParam) ([]models.OrderHistoryResponse, error) {

	var queryResult []models.OrderHistoryResult

	query := r.Database.WithContext(ctx).
		Table("orders AS o").
		Select(`
		o.id, 
		o.total_qty, 
		o.amount, 
		o.status, 
		o.payment_method, 
		o.shipping_address, 
		od.products, 
		od.order_history`).
		Joins("JOIN order_detail AS od ON od.id = o.order_detail_id").
		Where("o.user_id = ?", param.UserID)

	if param.Status > 0 {
		query.Where("o.status = ?", param.Status)
	}

	err := query.Order("o.id DESC").Scan(&queryResult).Error
	if err != nil {
		return nil, err
	}

	var results []models.OrderHistoryResponse
	for _, result := range queryResult {
		var products []models.CheckoutItem
		var history []models.StatusHistory

		err = json.Unmarshal([]byte(result.Products), &products)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(result.History), &history)
		if err != nil {
			return nil, err
		}

		results = append(results, models.OrderHistoryResponse{
			OrderID:         result.ID,
			TotalAmount:     result.Amount,
			TotalQty:        result.TotalQty,
			Status:          constant.OrderStatusTranslated[result.Status],
			PaymentMethod:   result.PaymentMethod,
			ShippingAddress: result.ShippingAddress,
			Products:        products,
			History:         history,
		})
	}

	return results, nil
}

package models

import "time"

type Order struct {
	ID              int64   `json:"id"`
	UserID          int64   `json:"user_id"`
	Amount          float64 `json:"amount"`
	TotalQty        int     `json:"total_qty"`
	OrderDetailID   int64   `json:"order_detail_id"`
	Status          int     `json:"status"`
	PaymentMethod   string  `json:"payment_method"`
	ShippingAddress string  `json:"shipping_address"`
}

type OrderDetail struct {
	ID           int64
	Products     string // stringfy json
	OrderHistory string // stringfy json
}

type CheckoutItem struct {
	ProductID int64   `json:"product_id"`
	Quantity  int64   `json:"quantity"`
	Price     float64 `json:"price"`
}

type CheckoutRequest struct {
	UserID           int64          `json:"user_id"`
	Items            []CheckoutItem `json:"items"`
	PaymentMethod    string         `json:"payment_method"`
	ShippingAddress  string         `json:"shipping_address"`
	IdempotencyToken string         `json:"idempotency_token"`
}

type OrderRequestLog struct {
	ID               int64     `json:"id"`
	IdempotencyToken string    `json:"idempotency_token"`
	CreateTime       time.Time `json:"create_time"`
}

type OrderHistoryParam struct {
	UserID int64
	Status int
}

type OrderHistoryResponse struct {
	OrderID         int64           `json:"order_id"`
	TotalAmount     float64         `json:"total_amount"`
	TotalQty        int             `json:"total_qty"`
	Status          string          `json:"status"`
	PaymentMethod   string          `json:"payment_method"`
	ShippingAddress string          `json:"shipping_address"`
	Products        []CheckoutItem  `json:"products"`
	History         []StatusHistory `json:"history"`
}

type StatusHistory struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

type OrderHistoryResult struct {
	ID              int64 `gorm:"column:id"`
	Amount          float64
	TotalQty        int
	Status          int
	PaymentMethod   string
	ShippingAddress string
	Products        string `gorm:"column:products"`
	History         string `gorm:"column:order_history"`
}

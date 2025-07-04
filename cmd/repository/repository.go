package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"net/http"
	"order_service/infra/log"
	"order_service/models"
)

type OrderRepository struct {
	Database    *gorm.DB
	Redis       *redis.Client
	ProductHost string
}

func NewOrderRepository(db *gorm.DB, redisClient *redis.Client, productHost string) *OrderRepository {
	return &OrderRepository{
		Database:    db,
		Redis:       redisClient,
		ProductHost: productHost,
	}
}

func (r *OrderRepository) GetProductInfo(ctx context.Context, productId int64) (models.Product, error) {
	var model models.Product

	url := fmt.Sprintf("%s/v1/product/%d", r.ProductHost, productId)
	log.Logger.Info(fmt.Sprintf("Requesting product info from %s", url))

	// Create the HTTP request with the provided context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productId": productId,
			"url":       url,
			"err":       err.Error(),
		}).Error("Failed to create HTTP request")
		return models.Product{}, err
	}

	// Perform the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productId": productId,
			"url":       url,
			"err":       err.Error(),
		}).Error("Failed to execute HTTP request")
		return models.Product{}, errors.New("failed to get product detail")
	}

	// Ensure the response body is closed once the function completes
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Logger.WithFields(logrus.Fields{
				"productId": productId,
				"err":       err.Error(),
			}).Error("Failed to close response body")
		}
	}(resp.Body)

	// Check if the status code is OK (200)
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return models.Product{}, nil
		}

		log.Logger.WithFields(logrus.Fields{
			"productId":  productId,
			"statusCode": resp.StatusCode,
			"url":        url,
		}).Error("Received non-200 HTTP status")
		return models.Product{}, fmt.Errorf("invalid response, status code: %d", resp.StatusCode)
	}

	// Decode the response body into the model
	var response models.GetProductInfo
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productId": productId,
			"err":       err.Error(),
		}).Error("Failed to decode response body")
		return models.Product{}, err
	}
	model = response.Product

	log.Logger.WithFields(logrus.Fields{
		"productId": productId,
		"product":   model,
	}).Info("Successfully retrieved product info")

	return model, nil
}

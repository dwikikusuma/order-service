package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"order_service/cmd/usecase"
	"order_service/infra/log"
	"order_service/infra/utils"
	"order_service/models"
	"strconv"
)

type OrderHandler struct {
	OrderUseCase usecase.OrderUseCase
}

func NewHandler(orderUseCase usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		OrderUseCase: orderUseCase,
	}
}

func (h *OrderHandler) CheckOutOrder(c *gin.Context) {
	var param models.CheckoutRequest
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	userId, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	if len(param.Items) == 0 || param.Items == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "item cannot be null",
		})
		return
	}

	param.UserID = int64(userId)
	orderId, err := h.OrderUseCase.CheckOutOrder(c.Request.Context(), &param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to checkout order",
		})

		log.Logger.WithFields(logrus.Fields{
			"param": param,
			"err":   err,
		}).Error("failed to checkout order")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ok":       true,
		"order_id": orderId,
	})

}

func (h *OrderHandler) GetOrderHistory(c *gin.Context) {
	userIdF, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}
	userId := int64(userIdF)

	statusStr := c.DefaultQuery("status", "0")
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid status",
		})
		return
	}

	param := &models.OrderHistoryParam{
		UserID: userId,
		Status: status,
	}

	history, err := h.OrderUseCase.GetOrderHistoryByUserId(c.Request.Context(), param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param":   param,
			"message": "error occurred on h.OrderUseCase.GetOrderHistoryByUserId(c.Request.Context(), param)",
			"error":   err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "invalid failed to get user history",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": history,
	})
}

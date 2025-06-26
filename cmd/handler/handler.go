package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"order_service/cmd/usecase"
	"order_service/infra/log"
	"order_service/infra/utils"
	"order_service/models"
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

	param.UserID = userId
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

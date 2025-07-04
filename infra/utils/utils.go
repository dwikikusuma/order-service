package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) (float64, error) {
	v, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("unauthorized")
	}
	id, ok := v.(float64)
	if !ok {
		return 0, errors.New("invalid user_id")
	}
	return id, nil
}

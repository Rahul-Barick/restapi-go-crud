package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetReferenceID(c *gin.Context) (string, error) {
	refID, exists := c.Get("referenceId")
	if !exists {
		return "", errors.New("missing idempotency key in referenceId header")
	}
	refIDStr, ok := refID.(string)
	if !ok || refIDStr == "" {
		return "", errors.New("invalid idempotency key format in referenceId header")
	}
	return refIDStr, nil
}

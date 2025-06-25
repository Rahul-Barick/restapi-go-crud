package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware to extract idempotency key (referenceId) from header and set in context
func IdempotencyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		referenceID := c.GetHeader("referenceId")
		if referenceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing idempotency key in referenceId header"})
			c.Abort()
			return
		}
		c.Set("referenceId", referenceID)
		c.Next()
	}
}

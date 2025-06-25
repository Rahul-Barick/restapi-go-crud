package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func MapValidationErrors(c *gin.Context, err error) bool {
	// Handle type conversion errors (e.g., string instead of int)
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		field := strings.ToLower(typeErr.Field)
		var message string

		switch field {
		case "account_id":
			message = "account_id must be an integer"
		case "initial_balance":
			message = "initial_balance must be a decimal number"
		default:
			message = field + " has an invalid type"
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return true
	}

	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		for _, e := range verrs {
			switch e.Field() {
			case "InitialBalance":
				if e.Tag() == "nonnegdecimal" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "initial_balance cannot be negative"})
					return true
				}
				if e.Tag() == "required" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "initial_balance is required"})
					return true
				}
			case "AccountID":
				if e.Tag() == "gt" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "account_id must be greater than 0"})
					return true
				}
				if e.Tag() == "required" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
					return true
				}
			}
		}
	}
	return false
}

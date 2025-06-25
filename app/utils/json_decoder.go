package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// StrictJSONBind is a helper function that binds and validates JSON requests strictly.
// It ensures the incoming JSON body does not contain unknown fields and properly matches the expected structure.
// On validation or decoding error, it returns a detailed HTTP 400 error message to the client.
// Returns true if the JSON is valid and successfully bound to the target object; otherwise returns false.
func StrictJSONBind(c *gin.Context, obj interface{}) bool {
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(obj); err != nil {
		switch e := err.(type) {
		case *json.SyntaxError:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Request body contains badly-formed JSON at position %d", e.Offset),
			})
		case *json.UnmarshalTypeError:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("%s must be a %s", e.Field, e.Type),
			})
		default:
			if strings.HasPrefix(err.Error(), "json: unknown field ") {
				field := strings.TrimPrefix(err.Error(), "json: unknown field ")
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("Request contains unknown field: %s", field),
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		}
		return false
	}
	return true
}

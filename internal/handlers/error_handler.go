package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				if !c.Writer.Written() {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "internal server error" + e.Error(),
					})
					return
				}
			}
		}
	}
}

func handleValidationErrors(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		errorMessage := make([]string, 0)
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errorMessage = append(errorMessage, fmt.Sprintf("%s is required", e.Field()))
			case "email":
				errorMessage = append(errorMessage, fmt.Sprintf("%s must be a valid email", e.Field()))
			case "min":
				errorMessage = append(errorMessage, fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param()))
			case "max":
				errorMessage = append(errorMessage, fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param()))
			default:
				errorMessage = append(errorMessage, fmt.Sprintf("%s is invalid", e.Field()))
			}
		}
		c.JSON(http.StatusBadRequest, ErrorResponse{Errors: errorMessage})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
}

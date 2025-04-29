package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mas963/go_rest_api/internal/services"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(http.StatusInternalServerError, services.ErrorResponse{
				Code:    services.ErrInternal.Code,
				Message: err.Error(),
			})
		}
	}
}

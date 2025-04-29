package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mas963/go_rest_api/internal/services"
)

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != requiredRole {
			c.JSON(http.StatusForbidden, services.ErrorResponse{
				Code:    services.ErrUnauthorized.Code,
				Message: "Insufficient permissions",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mas963/go_rest_api/internal/models"
	"github.com/mas963/go_rest_api/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, services.ErrorResponse{
			Code:    services.ErrValidation.Code,
			Message: "Invalid input: " + err.Error(),
		})
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		status := http.StatusUnauthorized
		if err == services.ErrNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, services.ErrorResponse{
			Code:    services.GetErrorCode(err),
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, services.ErrorResponse{
			Code:    services.ErrValidation.Code,
			Message: "Invalid input: " + err.Error(),
		})
		return
	}

	if err := h.authService.Register(req); err != nil {
		c.JSON(http.StatusInternalServerError, services.ErrorResponse{
			Code:    services.GetErrorCode(err),
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
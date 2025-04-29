package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mas963/go_rest_api/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, services.ErrorResponse{
			Code:    services.ErrValidation.Code,
			Message: "User ID is required",
		})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, services.ErrorResponse{
			Code:    services.GetErrorCode(err),
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

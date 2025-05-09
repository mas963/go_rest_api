package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mas963/go_rest_api/internal/models"
	"github.com/mas963/go_rest_api/internal/services"
)

type TokenResponse struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func Register(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var createDTO models.CreateUserDTO
		if err := c.ShouldBindJSON(&createDTO); err != nil {
			handleValidationErrors(c, err)
			return
		}

		user, err := userService.Create(createDTO)
		if err != nil {
			if err == services.ErrEmailAlreadyExists {
				c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
				return
			}
			c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

func Login(userService services.UserService, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginDTO models.LoginDTO
		if err := c.ShouldBindJSON(&loginDTO); err != nil {
			handleValidationErrors(c, err)
			return
		}

		user, err := userService.Authenticate(loginDTO.Email, loginDTO.Password)
		if err != nil {
			if err == services.ErrInvalidCredentials {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
				return
			}
			c.Error(err)
			return
		}

		expirationTime := time.Now().Add(24 * time.Hour)
		claims := jwt.MapClaims{
			"user_id": user.ID,
			"exp":     expirationTime.Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, TokenResponse{
			Token:   tokenString,
			Expires: expirationTime,
		})
	}
}

func GetUsers(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := userService.GetAll()
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func GetUser(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		user, err := userService.GetByID(uint(id))
		if err != nil {
			if err == services.ErrUserNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func UpdateUser(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		var updateDTO models.UpdateUserDTO
		if err := c.ShouldBindJSON(&updateDTO); err != nil {
			handleValidationErrors(c, err)
			return
		}

		claims, _ := c.Get("claims")
		if claims != nil {
			userID := claims.(jwt.MapClaims)["user_id"].(float64)
			if uint(userID) != uint(id) {
				c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized user"})
				return
			}
		}

		user, err := userService.Update(uint(id), updateDTO)
		if err != nil {
			if err == services.ErrUserNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			if err == services.ErrEmailAlreadyExists {
				c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
				return
			}
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func DeleteUser(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		claims, _ := c.Get("claims")
		if claims != nil {
			userID := claims.(jwt.MapClaims)["user_id"].(float64)
			if uint(userID) != uint(id) {
				c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized user"})
				return
			}
		}

		err = userService.Delete(uint(id))
		if err != nil {
			if err == services.ErrUserNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
	}
}

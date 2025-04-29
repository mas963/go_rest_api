package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mas963/go_rest_api/internal/config"
	"github.com/mas963/go_rest_api/internal/handlers"
	"github.com/mas963/go_rest_api/internal/middleware"
	"github.com/mas963/go_rest_api/internal/repositories"
	"github.com/mas963/go_rest_api/internal/services"
	"github.com/mas963/go_rest_api/internal/validators"
)

func main() {
	cfg := config.LoadConfig()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("strongpassword", validators.StrongPasswordValidator)
	}

	userRepo := repositories.NewUserRepository(cfg.DB)

	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)

	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.Default()
	r.Use(handlers.ErrorHandler()) // add global error handler

	// public routes
	r.POST("/login", authHandler.Login)
	r.POST("/register", authHandler.Register)

	// protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// user routes
		api.GET("/users/:id", userHandler.GetUser)

		// admin-only routes
		admin := api.Group("/admin").Use(middleware.RoleMiddleware("admin"))
		{
			admin.GET("/users", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Admin-only endpoint"})
			})
		}
	}

	r.Run(":8080")
}

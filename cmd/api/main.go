package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mas963/go_rest_api/internal/config"
	"github.com/mas963/go_rest_api/internal/handlers"
	"github.com/mas963/go_rest_api/internal/middleware"
	"github.com/mas963/go_rest_api/internal/repositories"
	"github.com/mas963/go_rest_api/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := config.SetupDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)

	userService := services.NewUserService(userRepo)

	router := gin.Default()

	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	v1 := router.Group("/api/v1")
	{
		v1.POST("/login", handlers.Login(userService, cfg.JWTSecret))
		v1.POST("/register", handlers.Register(userService))

		authorized := v1.Group("/")
		authorized.Use(authMiddleware.Authenticate())
		{
			authorized.GET("/users", handlers.GetUsers(userService))
			authorized.GET("/users/:id", handlers.GetUser(userService))
			authorized.PUT("/users/:id", handlers.UpdateUser(userService))
			authorized.DELETE("/users/:id", handlers.DeleteUser(userService))
		}
	}

	router.Use(handlers.ErrorHandler())

	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatal(err)
	}
}

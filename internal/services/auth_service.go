package services

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mas963/go_rest_api/internal/models"
	"github.com/mas963/go_rest_api/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: os.Getenv("JWT_SECRET"),
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) Login(email, password string) (string, error) {
	if email == "" || password == "" {
		return "", ErrValidation
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", ErrNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrUnauthorized
	}

	claims := Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) Register(req models.RegisterRequest) error {
	// password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Password = string(hashedPassword)

	return s.userRepo.Create(&models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     "user",
	})
}

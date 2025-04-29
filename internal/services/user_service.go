package services

import (
	"github.com/mas963/go_rest_api/internal/models"
	"github.com/mas963/go_rest_api/internal/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	if id == "" {
		return nil, ErrValidation
	}
	return s.userRepo.FindByID(id)
}
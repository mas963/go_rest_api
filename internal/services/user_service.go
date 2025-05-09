package services

import (
	"errors"

	"github.com/mas963/go_rest_api/internal/models"
	"github.com/mas963/go_rest_api/internal/repositories"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserService interface {
	Create(dto models.CreateUserDTO) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(id uint, dto models.UpdateUserDTO) (*models.User, error)
	Delete(id uint) error
	Authenticate(email, password string) (*models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *userService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Create(dto models.CreateUserDTO) (*models.User, error) {
	existingUser, err := s.userRepo.FindByEmail(dto.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	user := &models.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		Role:     "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *userService) GetByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *userService) GetAll() ([]models.User, error) {
	return s.userRepo.FindAll()
}

func (s *userService) Update(id uint, dto models.UpdateUserDTO) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if dto.Email != "" && dto.Email != user.Email {
		existingUser, err := s.userRepo.FindByEmail(dto.Email)
		if err != nil {
			return nil, err
		}
		if existingUser != nil {
			return nil, ErrEmailAlreadyExists
		}
		user.Email = dto.Email
	}

	if dto.Name != "" {
		user.Name = dto.Name
	}

	if dto.Password != "" {
		user.Password = dto.Password
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}


func (s *userService) Delete(id uint) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	return s.userRepo.Delete(id)
}

func (s *userService) Authenticate(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil || !user.CheckPassword(password) {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}


package usecase

import (
	"fmt"

	"github.com/ucups/go-user-service/internal/domain"
	"github.com/ucups/go-user-service/internal/repository"
)

// UserUseCase handles user business logic
type UserUseCase struct {
	repo repository.UserRepository
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

// CreateUser creates a new user
func (uc *UserUseCase) CreateUser(name string) (*domain.User, error) {
	// Create and validate user
	user, err := domain.NewUser(name)
	if err != nil {
		return nil, err
	}

	// Persist user
	if err := uc.repo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (uc *UserUseCase) GetUserByID(id int64) (*domain.User, error) {
	user, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetAllUsers retrieves all users with pagination
func (uc *UserUseCase) GetAllUsers(pageNum, pageSize int) ([]*domain.User, error) {
	// Set defaults
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	filter := domain.UserFilter{
		PageNum:  pageNum,
		PageSize: pageSize,
	}

	users, err := uc.repo.GetAll(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

package repository

import "github.com/ucups/go-user-service/internal/domain"

// UserRepository defines the interface for user data persistence
type UserRepository interface {
	// Create adds a new user to the repository
	Create(user *domain.User) error

	// GetByID retrieves a user by ID
	GetByID(id int64) (*domain.User, error)

	// GetAll retrieves all users with pagination
	GetAll(filter domain.UserFilter) ([]*domain.User, error)

	// Close closes the repository connection
	Close() error
}

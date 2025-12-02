package domain

import (
	"errors"
	"strings"
	"time"
)

// User represents a user entity
type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// NewUser creates a new user with validation
func NewUser(name string) (*User, error) {
	if err := ValidateName(name); err != nil {
		return nil, err
	}

	now := time.Now().UnixMicro()
	return &User{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// ValidateName validates the user name
func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("name is required")
	}
	return nil
}

// UserFilter represents filter options for querying users
type UserFilter struct {
	PageNum  int
	PageSize int
}

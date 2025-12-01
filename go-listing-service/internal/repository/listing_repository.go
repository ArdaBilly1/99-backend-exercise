package repository

import (
	"context"

	"github.com/ucups/go-listing-service/internal/domain"
)

// ListingRepository defines the interface for listing data persistence
type ListingRepository interface {
	Create(ctx context.Context, listing *domain.Listing) error
	FindAll(ctx context.Context, filter domain.ListingFilter) ([]*domain.Listing, error)
	Close() error
}

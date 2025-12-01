package usecase

import (
	"context"
	"fmt"

	"github.com/ucups/go-listing-service/internal/domain"
	"github.com/ucups/go-listing-service/internal/repository"
)

// ListingUseCase defines the interface for listing business logic
type ListingUseCase interface {
	CreateListing(ctx context.Context, userID int64, listingType string, price int64) (*domain.Listing, error)
	GetListings(ctx context.Context, userID *int64, pageNum, pageSize int) ([]*domain.Listing, error)
}

type listingUseCase struct {
	repo repository.ListingRepository
}

// NewListingUseCase creates a new listing use case
func NewListingUseCase(repo repository.ListingRepository) ListingUseCase {
	return &listingUseCase{
		repo: repo,
	}
}

// CreateListing handles the business logic for creating a listing
func (u *listingUseCase) CreateListing(ctx context.Context, userID int64, listingType string, price int64) (*domain.Listing, error) {
	// Convert string to domain type
	lt := domain.ListingType(listingType)

	// Create listing with domain validation
	listing, err := domain.NewListing(userID, lt, price)
	if err != nil {
		return nil, err
	}

	// Persist to repository
	if err := u.repo.Create(ctx, listing); err != nil {
		return nil, fmt.Errorf("failed to save listing: %w", err)
	}

	return listing, nil
}

// GetListings handles the business logic for retrieving listings
func (u *listingUseCase) GetListings(ctx context.Context, userID *int64, pageNum, pageSize int) ([]*domain.Listing, error) {
	// Set default pagination if not provided
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	filter := domain.ListingFilter{
		UserID:   userID,
		PageNum:  pageNum,
		PageSize: pageSize,
	}

	listings, err := u.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve listings: %w", err)
	}

	return listings, nil
}

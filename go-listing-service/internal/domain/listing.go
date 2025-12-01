package domain

import (
	"errors"
	"time"
)

// ListingType represents the type of listing
type ListingType string

const (
	ListingTypeRent ListingType = "rent"
	ListingTypeSale ListingType = "sale"
)

// Listing represents a property listing entity
type Listing struct {
	ID          int64       `json:"id"`
	UserID      int64       `json:"user_id"`
	ListingType ListingType `json:"listing_type"`
	Price       int64       `json:"price"`
	CreatedAt   int64       `json:"created_at"`
	UpdatedAt   int64       `json:"updated_at"`
}

// NewListing creates a new listing with validation
func NewListing(userID int64, listingType ListingType, price int64) (*Listing, error) {
	if err := ValidateListingType(listingType); err != nil {
		return nil, err
	}
	if err := ValidatePrice(price); err != nil {
		return nil, err
	}

	now := time.Now().UnixMicro()
	return &Listing{
		UserID:      userID,
		ListingType: listingType,
		Price:       price,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// ValidateListingType validates the listing type
func ValidateListingType(listingType ListingType) error {
	if listingType != ListingTypeRent && listingType != ListingTypeSale {
		return errors.New("listing_type must be either 'rent' or 'sale'")
	}
	return nil
}

// ValidatePrice validates the price
func ValidatePrice(price int64) error {
	if price <= 0 {
		return errors.New("price must be greater than 0")
	}
	return nil
}

// ListingFilter represents filter options for querying listings
type ListingFilter struct {
	UserID   *int64
	PageNum  int
	PageSize int
}

package model

// User represents user data from user service
type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// Listing represents listing data from listing service
type Listing struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	ListingType string `json:"listing_type"`
	Price       int64  `json:"price"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

// EnrichedListing represents a listing with embedded user information
type EnrichedListing struct {
	ID          int64  `json:"id"`
	ListingType string `json:"listing_type"`
	Price       int64  `json:"price"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	User        User   `json:"user"`
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Name string `json:"name"`
}

// CreateListingRequest represents the request to create a listing
type CreateListingRequest struct {
	UserID      int64  `json:"user_id"`
	ListingType string `json:"listing_type"`
	Price       int64  `json:"price"`
}

// ServiceResponse represents a generic service response
type ServiceResponse struct {
	Result bool                   `json:"result"`
	Data   map[string]interface{} `json:"data,omitempty"`
	Errors []string               `json:"errors,omitempty"`
}

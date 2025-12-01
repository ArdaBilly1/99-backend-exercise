package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ucups/go-listing-service/internal/domain"
)

type listingRepository struct {
	db *sql.DB
}

// NewListingRepository creates a new SQLite listing repository
func NewListingRepository(dbPath string) (*listingRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create table if not exists
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS listings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		listing_type TEXT NOT NULL,
		price INTEGER NOT NULL,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL
	)`

	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &listingRepository{db: db}, nil
}

// Create inserts a new listing into the database
func (r *listingRepository) Create(ctx context.Context, listing *domain.Listing) error {
	query := `
		INSERT INTO listings (user_id, listing_type, price, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		listing.UserID,
		listing.ListingType,
		listing.Price,
		listing.CreatedAt,
		listing.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create listing: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	listing.ID = id
	return nil
}

// FindAll retrieves listings with optional filtering and pagination
func (r *listingRepository) FindAll(ctx context.Context, filter domain.ListingFilter) ([]*domain.Listing, error) {
	query := "SELECT id, user_id, listing_type, price, created_at, updated_at FROM listings"
	args := []interface{}{}

	// Add user_id filter if provided
	if filter.UserID != nil {
		query += " WHERE user_id = ?"
		args = append(args, *filter.UserID)
	}

	// Add pagination
	query += " ORDER BY id DESC LIMIT ? OFFSET ?"
	offset := (filter.PageNum - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query listings: %w", err)
	}
	defer rows.Close()

	var listings []*domain.Listing
	for rows.Next() {
		listing := &domain.Listing{}
		var listingType string

		err := rows.Scan(
			&listing.ID,
			&listing.UserID,
			&listingType,
			&listing.Price,
			&listing.CreatedAt,
			&listing.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan listing: %w", err)
		}

		listing.ListingType = domain.ListingType(listingType)
		listings = append(listings, listing)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return listings, nil
}

// Close closes the database connection
func (r *listingRepository) Close() error {
	return r.db.Close()
}

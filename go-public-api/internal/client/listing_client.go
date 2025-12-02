package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ucups/go-public-api/internal/model"
)

// ListingClient handles communication with listing service
type ListingClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewListingClient creates a new listing service client
func NewListingClient(baseURL string) *ListingClient {
	return &ListingClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// GetListings retrieves listings with optional user_id filter
func (c *ListingClient) GetListings(pageNum, pageSize int, userID *int64) ([]model.Listing, error) {
	params := url.Values{}
	params.Add("page_num", strconv.Itoa(pageNum))
	params.Add("page_size", strconv.Itoa(pageSize))
	if userID != nil {
		params.Add("user_id", strconv.FormatInt(*userID, 10))
	}

	apiURL := fmt.Sprintf("%s/listings?%s", c.baseURL, params.Encode())
	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get listings: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("listing service returned status %d: %s", resp.StatusCode, string(body))
	}

	var serviceResp model.ServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&serviceResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !serviceResp.Result {
		return nil, fmt.Errorf("listing service error: %v", serviceResp.Errors)
	}

	// Extract listings from data
	listingsData, ok := serviceResp.Data["listings"]
	if !ok {
		return []model.Listing{}, nil
	}

	// Convert to JSON and back to proper type
	jsonData, err := json.Marshal(listingsData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal listings: %w", err)
	}

	var listings []model.Listing
	if err := json.Unmarshal(jsonData, &listings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal listings: %w", err)
	}

	return listings, nil
}

// CreateListing creates a new listing
func (c *ListingClient) CreateListing(userID int64, listingType string, price int64) (*model.Listing, error) {
	data := url.Values{}
	data.Set("user_id", strconv.FormatInt(userID, 10))
	data.Set("listing_type", listingType)
	data.Set("price", strconv.FormatInt(price, 10))

	apiURL := fmt.Sprintf("%s/listings", c.baseURL)
	resp, err := c.httpClient.PostForm(apiURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create listing: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("listing service returned status %d: %s", resp.StatusCode, string(body))
	}

	var serviceResp model.ServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&serviceResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !serviceResp.Result {
		return nil, fmt.Errorf("listing service error: %v", serviceResp.Errors)
	}

	// Extract listing from data
	listingData, ok := serviceResp.Data["listing"]
	if !ok {
		return nil, fmt.Errorf("no listing in response")
	}

	// Convert to JSON and back to proper type
	jsonData, err := json.Marshal(listingData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal listing: %w", err)
	}

	var listing model.Listing
	if err := json.Unmarshal(jsonData, &listing); err != nil {
		return nil, fmt.Errorf("failed to unmarshal listing: %w", err)
	}

	return &listing, nil
}

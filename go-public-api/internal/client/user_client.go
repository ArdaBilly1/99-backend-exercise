package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ucups/go-public-api/internal/model"
)

// UserClient handles communication with user service
type UserClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewUserClient creates a new user service client
func NewUserClient(baseURL string) *UserClient {
	return &UserClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// GetUser retrieves a user by ID
func (c *UserClient) GetUser(userID int64) (*model.User, error) {
	apiURL := fmt.Sprintf("%s/users/%d", c.baseURL, userID)
	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user service returned status %d: %s", resp.StatusCode, string(body))
	}

	var serviceResp model.ServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&serviceResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !serviceResp.Result {
		return nil, fmt.Errorf("user service error: %v", serviceResp.Errors)
	}

	// Extract user from data
	userData, ok := serviceResp.Data["user"]
	if !ok {
		return nil, fmt.Errorf("no user in response")
	}

	// Convert to JSON and back to proper type
	jsonData, err := json.Marshal(userData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user: %w", err)
	}

	var user model.User
	if err := json.Unmarshal(jsonData, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

// CreateUser creates a new user
func (c *UserClient) CreateUser(name string) (*model.User, error) {
	data := url.Values{}
	data.Set("name", name)

	apiURL := fmt.Sprintf("%s/users", c.baseURL)
	resp, err := c.httpClient.PostForm(apiURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user service returned status %d: %s", resp.StatusCode, string(body))
	}

	var serviceResp model.ServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&serviceResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !serviceResp.Result {
		return nil, fmt.Errorf("user service error: %v", serviceResp.Errors)
	}

	// Extract user from data
	userData, ok := serviceResp.Data["user"]
	if !ok {
		return nil, fmt.Errorf("no user in response")
	}

	// Convert to JSON and back to proper type
	jsonData, err := json.Marshal(userData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user: %w", err)
	}

	var user model.User
	if err := json.Unmarshal(jsonData, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

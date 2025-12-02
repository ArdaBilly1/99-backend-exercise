package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ucups/go-public-api/internal/client"
	"github.com/ucups/go-public-api/internal/model"
)

// PublicHandler handles public API requests
type PublicHandler struct {
	listingClient *client.ListingClient
	userClient    *client.UserClient
}

// NewPublicHandler creates a new public API handler
func NewPublicHandler(listingClient *client.ListingClient, userClient *client.UserClient) *PublicHandler {
	return &PublicHandler{
		listingClient: listingClient,
		userClient:    userClient,
	}
}

// GetListings handles GET /public-api/listings
func (h *PublicHandler) GetListings(w http.ResponseWriter, r *http.Request) {
	// Parse pagination params
	pageNumStr := r.URL.Query().Get("page_num")
	pageSizeStr := r.URL.Query().Get("page_size")
	userIDStr := r.URL.Query().Get("user_id")

	pageNum := 1
	pageSize := 10
	var userID *int64

	if pageNumStr != "" {
		if val, err := strconv.Atoi(pageNumStr); err == nil {
			pageNum = val
		} else {
			WriteError(w, http.StatusBadRequest, "invalid page_num")
			return
		}
	}

	if pageSizeStr != "" {
		if val, err := strconv.Atoi(pageSizeStr); err == nil {
			pageSize = val
		} else {
			WriteError(w, http.StatusBadRequest, "invalid page_size")
			return
		}
	}

	if userIDStr != "" {
		if val, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			userID = &val
		} else {
			WriteError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
	}

	// Get listings from listing service
	listings, err := h.listingClient.GetListings(pageNum, pageSize, userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Enrich each listing with user data
	enrichedListings := make([]model.EnrichedListing, 0, len(listings))
	for _, listing := range listings {
		user, err := h.userClient.GetUser(listing.UserID)
		if err != nil {
			// If user not found, skip this listing or use placeholder
			// For now, we'll return an error
			WriteError(w, http.StatusInternalServerError, "failed to get user data: "+err.Error())
			return
		}

		enrichedListing := model.EnrichedListing{
			ID:          listing.ID,
			ListingType: listing.ListingType,
			Price:       listing.Price,
			CreatedAt:   listing.CreatedAt,
			UpdatedAt:   listing.UpdatedAt,
			User:        *user,
		}
		enrichedListings = append(enrichedListings, enrichedListing)
	}

	// Return response
	WriteSuccess(w, map[string]interface{}{
		"result":   true,
		"listings": enrichedListings,
	})
}

// CreateUser handles POST /public-api/users
func (h *PublicHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Create user via user service
	user, err := h.userClient.CreateUser(req.Name)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return response
	WriteSuccess(w, map[string]interface{}{
		"user": user,
	})
}

// CreateListing handles POST /public-api/listings
func (h *PublicHandler) CreateListing(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	var req model.CreateListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Create listing via listing service
	listing, err := h.listingClient.CreateListing(req.UserID, req.ListingType, req.Price)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return response
	WriteSuccess(w, map[string]interface{}{
		"listing": listing,
	})
}

// Ping handles GET /public-api/ping
func (h *PublicHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong!"))
}

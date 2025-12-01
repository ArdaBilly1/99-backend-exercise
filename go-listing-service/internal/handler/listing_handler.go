package handler

import (
	"net/http"
	"strconv"

	"github.com/ucups/go-listing-service/internal/usecase"
)

type ListingHandler struct {
	useCase usecase.ListingUseCase
}

// NewListingHandler creates a new listing handler
func NewListingHandler(useCase usecase.ListingUseCase) *ListingHandler {
	return &ListingHandler{
		useCase: useCase,
	}
}

// CreateListing handles POST /listings
func (h *ListingHandler) CreateListing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid form data")
		return
	}

	// Validate and parse user_id
	userIDStr := r.FormValue("user_id")
	if userIDStr == "" {
		WriteError(w, http.StatusBadRequest, "user_id is required")
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "user_id must be a valid integer")
		return
	}

	// Validate listing_type
	listingType := r.FormValue("listing_type")
	if listingType == "" {
		WriteError(w, http.StatusBadRequest, "listing_type is required")
		return
	}

	// Validate and parse price
	priceStr := r.FormValue("price")
	if priceStr == "" {
		WriteError(w, http.StatusBadRequest, "price is required")
		return
	}
	price, err := strconv.ParseInt(priceStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "price must be a valid integer")
		return
	}

	// Create listing via use case
	listing, err := h.useCase.CreateListing(r.Context(), userID, listingType, price)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteSuccess(w, listing)
}

// GetListings handles GET /listings
func (h *ListingHandler) GetListings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Parse optional user_id filter
	var userID *int64
	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		id, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "user_id must be a valid integer")
			return
		}
		userID = &id
	}

	// Parse pagination parameters
	pageNum := 1
	if pageNumStr := r.URL.Query().Get("page_num"); pageNumStr != "" {
		num, err := strconv.Atoi(pageNumStr)
		if err != nil || num < 1 {
			WriteError(w, http.StatusBadRequest, "page_num must be a positive integer")
			return
		}
		pageNum = num
	}

	pageSize := 10
	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		size, err := strconv.Atoi(pageSizeStr)
		if err != nil || size < 1 {
			WriteError(w, http.StatusBadRequest, "page_size must be a positive integer")
			return
		}
		pageSize = size
	}

	// Get listings via use case
	listings, err := h.useCase.GetListings(r.Context(), userID, pageNum, pageSize)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteSuccess(w, listings)
}

// Ping handles GET /listings/ping
func (h *ListingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	WriteSuccess(w, map[string]string{"status": "ok"})
}

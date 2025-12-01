package handler

import (
	"net/http"

	"github.com/ucups/go-listing-service/internal/usecase"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(listingUseCase usecase.ListingUseCase) *http.ServeMux {
	mux := http.NewServeMux()

	// Initialize handler
	listingHandler := NewListingHandler(listingUseCase)

	// Register routes
	mux.HandleFunc("/listings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			listingHandler.CreateListing(w, r)
		case http.MethodGet:
			listingHandler.GetListings(w, r)
		default:
			WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	mux.HandleFunc("/listings/ping", listingHandler.Ping)

	return mux
}

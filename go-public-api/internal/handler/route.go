package handler

import (
	"github.com/gorilla/mux"
	"github.com/ucups/go-public-api/internal/client"
)

// SetupRoutes configures all HTTP routes
func SetupRoutes(listingClient *client.ListingClient, userClient *client.UserClient) *mux.Router {
	handler := NewPublicHandler(listingClient, userClient)

	router := mux.NewRouter()

	// Public API routes
	router.HandleFunc("/public-api/ping", handler.Ping).Methods("GET")
	router.HandleFunc("/public-api/listings", handler.GetListings).Methods("GET")
	router.HandleFunc("/public-api/listings", handler.CreateListing).Methods("POST")
	router.HandleFunc("/public-api/users", handler.CreateUser).Methods("POST")

	return router
}

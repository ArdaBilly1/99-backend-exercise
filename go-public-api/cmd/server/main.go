package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ucups/go-public-api/internal/client"
	"github.com/ucups/go-public-api/internal/config"
	"github.com/ucups/go-public-api/internal/handler"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize service clients
	listingClient := client.NewListingClient(cfg.Services.ListingServiceURL)
	userClient := client.NewUserClient(cfg.Services.UserServiceURL)

	// Setup routes
	mux := handler.SetupRoutes(listingClient, userClient)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	if cfg.Server.Debug {
		log.Printf("Starting public API in DEBUG mode on %s", addr)
		log.Printf("Listing Service: %s", cfg.Services.ListingServiceURL)
		log.Printf("User Service: %s", cfg.Services.UserServiceURL)
	} else {
		log.Printf("Starting public API on %s", addr)
	}

	// Graceful shutdown
	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

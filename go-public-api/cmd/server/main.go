package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ucups/go-public-api/internal/client"
	"github.com/ucups/go-public-api/internal/handler"
)

func main() {
	// Parse command-line flags
	port := flag.Int("port", 8000, "server port")
	debug := flag.Bool("debug", true, "debug mode")
	listingServiceURL := flag.String("listing-service", "http://localhost:6000", "listing service URL")
	userServiceURL := flag.String("user-service", "http://localhost:7000", "user service URL")
	flag.Parse()

	// Initialize service clients
	listingClient := client.NewListingClient(*listingServiceURL)
	userClient := client.NewUserClient(*userServiceURL)

	// Setup routes
	mux := handler.SetupRoutes(listingClient, userClient)

	// Start server
	addr := fmt.Sprintf(":%d", *port)
	if *debug {
		log.Printf("Starting public API in DEBUG mode on %s", addr)
		log.Printf("Listing Service: %s", *listingServiceURL)
		log.Printf("User Service: %s", *userServiceURL)
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

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ucups/go-listing-service/internal/handler"
	"github.com/ucups/go-listing-service/internal/repository/sqlite"
	"github.com/ucups/go-listing-service/internal/usecase"
)

func main() {
	// Parse command-line flags
	port := flag.Int("port", 6000, "server port")
	debug := flag.Bool("debug", true, "debug mode")
	flag.Parse()

	// Initialize repository layer
	repo, err := sqlite.NewListingRepository("listings.db")
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	// Initialize use case layer (dependency injection)
	listingUseCase := usecase.NewListingUseCase(repo)

	// Setup routes
	mux := handler.SetupRoutes(listingUseCase)

	// Start server
	addr := fmt.Sprintf(":%d", *port)
	if *debug {
		log.Printf("Starting server in DEBUG mode on %s", addr)
	} else {
		log.Printf("Starting server on %s", addr)
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

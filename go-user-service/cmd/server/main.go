package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ucups/go-user-service/internal/config"
	"github.com/ucups/go-user-service/internal/handler"
	"github.com/ucups/go-user-service/internal/repository/sqlite"
	"github.com/ucups/go-user-service/internal/usecase"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize repository layer
	repo, err := sqlite.NewUserRepository(cfg.DB.Path)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	// Initialize use case layer (dependency injection)
	userUseCase := usecase.NewUserUseCase(repo)

	// Setup routes
	mux := handler.SetupRoutes(userUseCase)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	if cfg.Server.Debug {
		log.Printf("Starting user service in DEBUG mode on %s", addr)
		log.Printf("Database path: %s", cfg.DB.Path)
		log.Printf("Listing Service URL: %s", cfg.Services.ListingServiceURL)
	} else {
		log.Printf("Starting user service on %s", addr)
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

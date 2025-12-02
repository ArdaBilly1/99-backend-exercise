package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ucups/go-user-service/internal/handler"
	"github.com/ucups/go-user-service/internal/repository/sqlite"
	"github.com/ucups/go-user-service/internal/usecase"
)

func main() {
	// Parse command-line flags
	port := flag.Int("port", 7000, "server port")
	debug := flag.Bool("debug", true, "debug mode")
	flag.Parse()

	// Initialize repository layer
	repo, err := sqlite.NewUserRepository("users.db")
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	// Initialize use case layer (dependency injection)
	userUseCase := usecase.NewUserUseCase(repo)

	// Setup routes
	mux := handler.SetupRoutes(userUseCase)

	// Start server
	addr := fmt.Sprintf(":%d", *port)
	if *debug {
		log.Printf("Starting user service in DEBUG mode on %s", addr)
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

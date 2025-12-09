package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	DB       DBConfig
	Services ServicesConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port  int
	Debug bool
}

// DBConfig holds database configuration
type DBConfig struct {
	Path string
}

// ServicesConfig holds external service URLs
type ServicesConfig struct {
	ListingServiceURL string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{}

	// Server configuration
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "7000" // Default port
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	config.Server.Port = port

	debugStr := os.Getenv("DEBUG_MODE")
	if debugStr == "" {
		debugStr = "true" // Default debug mode
	}
	config.Server.Debug = debugStr == "true"

	// Database configuration
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "users.db" // Default db path
	}
	config.DB.Path = dbPath

	// Services configuration
	config.Services.ListingServiceURL = os.Getenv("LISTING_SERVICE_URL")
	if config.Services.ListingServiceURL == "" {
		config.Services.ListingServiceURL = "http://localhost:6000" // Default listing service URL
	}

	return config, nil
}

# [config.go](config.go)

// Package config handles application configuration loading and validation.
// It provides functionality to load configuration from environment variables,
// with sensible defaults for development environments.


```go
package config

import (
    "log"
    "os"

    // Optional: helpful for local development. Loads .env
    "github.com/joho/godotenv"
)

// Config holds all application configuration settings loaded from environment variables.
type Config struct {
    // ServerPort is the HTTP port the server will listen on
    ServerPort string
    // ServerHost is the hostname for the server
    ServerHost string
    // DatabaseURL is the connection string for the PostgreSQL database
    DatabaseURL string
    // JWTSecret is the secret key used for JWT token signing and validation
    JWTSecret string
    // RedisAddr is the address of the Redis server for async job processing
    RedisAddr string
    // ExportDir is the directory where export files will be stored
    ExportDir string
    // Add more config fields as needed
}

// LoadConfig loads configuration from environment variables and returns a Config struct.
// If required environment variables are missing, it returns an error.
func LoadConfig() (*Config, error) {
    // Load .env file if it exists (useful for local development)
    // In production, environment variables are typically set directly.
    err := godotenv.Load()
    if err != nil && !os.IsNotExist(err) {
        log.Printf("Warning: Error loading .env file: %v\n", err)
    }

    cfg := &Config{
        ServerHost:  os.Getenv("HOST"),
        ServerPort:  os.Getenv("PORT"),
        DatabaseURL: os.Getenv("DATABASE_URL"),
        JWTSecret:   os.Getenv("JWT_SECRET"),
        RedisAddr:   os.Getenv("REDIS_ADDR"),
        ExportDir:   os.Getenv("EXPORT_DIR"),
        // Load other config...
    }

    // Basic validation (can be expanded)
    if cfg.ServerPort == "" {
        log.Println("Warning: PORT not set, using default 8080")
        cfg.ServerPort = "8080" // Default port
    }

    // Set default host if not provided
    if cfg.ServerHost == "" {
        log.Println("Warning: HOST not set, using default localhost")
        cfg.ServerHost = "localhost" // Default host
    }

    // Check if database URL is set, if not, return an error
    if cfg.DatabaseURL == "" {
        log.Fatal("DATABASE_URL is required")
    }

    // Check if JWT secret is set
    if cfg.JWTSecret == "" {
        cfg.JWTSecret = "development-jwt-secret" // Default for development
        log.Println("Warning: Using default JWT secret. Set JWT_SECRET for production.")
    }

    // Check if Redis address is set
    if cfg.RedisAddr == "" {
        cfg.RedisAddr = "localhost:6379" // Default Redis address
        log.Println("Warning: Using default Redis address. Set REDIS_ADDR for production.")
    }

    // Check if export directory is set
    if cfg.ExportDir == "" {
        cfg.ExportDir = "./exports" // Default export directory
        log.Println("Warning: Using default export directory. Set EXPORT_DIR for production.")
    }

    // Create export directory if it doesn't exist
    if err := os.MkdirAll(cfg.ExportDir, 0755); err != nil {
        log.Printf("Warning: Failed to create export directory %s: %v", cfg.ExportDir, err)
    }

    return cfg, nil
}
```

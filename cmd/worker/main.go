// Package main implements the ChatLogger API worker process.
// This worker handles asynchronous background jobs like chat exports,
// connecting to Redis for job processing and PostgreSQL for data access.
// It uses the asynq library for job queue management with graceful shutdown
// handling for production reliability.
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/kjanat/chatlogger-api-go/internal/jobs"
	"github.com/kjanat/chatlogger-api-go/internal/repository"
	"github.com/kjanat/chatlogger-api-go/internal/service"
)

// CustomLogger wraps the standard logger to implement asynq.Logger.
type CustomLogger struct {
	*log.Logger
}

// Debug logs debug messages.
func (l *CustomLogger) Debug(args ...interface{}) {
	l.Println(args...)
}

// Info logs info messages.
func (l *CustomLogger) Info(args ...interface{}) {
	l.Println(args...)
}

// Warn logs warning messages.
func (l *CustomLogger) Warn(args ...interface{}) {
	l.Println(args...)
}

// Error logs error messages.
func (l *CustomLogger) Error(args ...interface{}) {
	l.Println(args...)
}

// Fatal logs fatal messages.
func (l *CustomLogger) Fatal(args ...interface{}) {
	l.Println(args...)
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get Redis connection info
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	// Get export directory
	exportDir := os.Getenv("EXPORT_DIR")
	if exportDir == "" {
		exportDir = "./exports"
	}

	// Create exports directory if it doesn't exist
	if err := os.MkdirAll(exportDir, 0o750); err != nil {
		log.Fatalf("Failed to create export directory: %v", err)
	}

	// Get database connection info
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("DATABASE_URL environment variable is required")
	}

	// Initialize database connection without migrations
	dbOptions := repository.DefaultDatabaseOptions()
	dbOptions.RunMigrations = false // Worker should not run migrations

	db, err := repository.NewDatabaseWithOptions(dbURL, dbOptions)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Ensure the database connection is closed when the program exits
	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}

	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed successfully")
		}
	}()

	// Set up repositories
	exportRepo := repository.NewExportRepository(db.DB)
	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)

	// Set up services
	chatService := service.NewChatService(chatRepo)
	messageService := service.NewMessageService(messageRepo)

	// Create export processor
	processor := jobs.NewExportProcessor(
		exportRepo,
		chatService,
		messageService,
		exportDir,
	)

	// Create a custom logger that implements the asynq.Logger interface
	customLogger := &CustomLogger{
		Logger: log.New(os.Stdout, "asynq: ", log.LstdFlags),
	}

	// Set up Asynq server with shorter shutdown timeout
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 5, // Number of concurrent workers
			Queues: map[string]int{
				"exports": 5, // Process export queue with weight 5
				"default": 1, // Process default queue with weight 1
			},
			Logger:          customLogger,
			ShutdownTimeout: 5 * time.Second, // Force shutdown after 5 seconds
		},
	)

	// Register task handlers
	mux := asynq.NewServeMux()
	mux.HandleFunc(jobs.TypeExportProcess, processor.ProcessExport)

	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Create a done channel to indicate when shutdown is complete
	done := make(chan struct{})

	go func() {
		<-c
		log.Println("Shutting down worker...")

		// Set a timeout for shutdown to ensure we don't hang
		shutdownComplete := make(chan struct{})
		go func() {
			srv.Shutdown()
			close(shutdownComplete)
		}()

		// Wait for shutdown to complete or timeout
		select {
		case <-shutdownComplete:
			log.Println("Server shutdown completed normally")
		case <-time.After(8 * time.Second):
			log.Println("Server shutdown timed out - forcing exit")
		}

		// Signal that we're done
		close(done)

		// Force exit after a brief pause to allow logging to complete
		time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	}()

	// Start the worker
	log.Println("Starting export worker...")
	if err := srv.Run(mux); err != nil {
		log.Printf("Could not run server: %v", err)
		return
	}

	// Wait for shutdown to complete
	<-done
	log.Println("Worker shutdown complete")
}

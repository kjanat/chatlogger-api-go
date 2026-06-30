package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseOptions contains configuration options for database connections.
type DatabaseOptions struct {
	// RunMigrations determines whether to run database migrations on connect
	RunMigrations bool
}

// DefaultDatabaseOptions returns database options with sensible defaults.
func DefaultDatabaseOptions() *DatabaseOptions {
	return &DatabaseOptions{
		RunMigrations: true, // Default to running migrations
	}
}

// Database represents a database connection.
type Database struct {
	*gorm.DB
}

// NewDatabase creates a new database connection with default options.
func NewDatabase(dsn string) (*Database, error) {
	return NewDatabaseWithOptions(dsn, DefaultDatabaseOptions())
}

// NewDatabaseWithOptions creates a new database connection with the specified options.
func NewDatabaseWithOptions(dsn string, options *DatabaseOptions) (*Database, error) {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Only run migrations if enabled in options
	if options.RunMigrations {
		if err := runMigrations(db); err != nil {
			return nil, err
		}
	} else {
		log.Println("Database connected successfully (migrations skipped)")
	}

	return &Database{DB: db}, nil
}

// runMigrations runs the database schema migrations.
func runMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Begin a transaction to group all migration operations
	err := db.Transaction(func(tx *gorm.DB) error {
		// Auto migrate all models in a single transaction
		if err := tx.AutoMigrate(
			&domain.Organization{},
			&domain.APIKey{},
			&domain.User{},
			&domain.Chat{},
			&domain.Message{},
			&domain.Export{},
		); err != nil {
			return fmt.Errorf("failed to migrate database: %w", err)
		}

		log.Println("Database migrations completed successfully")
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to run database migration transaction: %w", err)
	}
	return nil
}

// Close closes the database connection.
func (db *Database) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB for closing: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return nil
}

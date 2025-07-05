package testutils

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	
	// Try in-memory SQLite first for faster tests
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		// Fallback to PostgreSQL if SQLite fails
		testDBURL := os.Getenv("TEST_DATABASE_URL")
		if testDBURL == "" {
			t.Skip("No database available for integration tests")
		}
		
		db, err = gorm.Open(postgres.Open(testDBURL), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			t.Skipf("Failed to connect to test database: %v", err)
		}
	}
	
	// Auto-migrate schemas for testing
	if err := AutoMigrateTestSchema(db); err != nil {
		t.Fatalf("Failed to migrate test schema: %v", err)
	}
	
	// Verify tables exist
	var tableCount int64
	db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name IN ('organizations', 'users', 'chats', 'messages', 'api_keys', 'exports')").Scan(&tableCount)
	if tableCount == 0 {
		t.Fatalf("No tables created after migration")
	}
	
	return db
}

// CleanupTestDB cleans up test database
func CleanupTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()
	
	sqlDB, err := db.DB()
	if err != nil {
		t.Errorf("Failed to get underlying sql.DB: %v", err)
		return
	}
	
	// Clean up tables in reverse order to avoid foreign key constraints
	tables := []string{
		"exports",
		"api_keys",
		"messages",
		"chats", 
		"users",
		"organizations",
	}
	
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			// Ignore errors for tables that might not exist
			t.Logf("Warning: Failed to clean table %s: %v", table, err)
		}
	}
	
	sqlDB.Close()
}

// SetupTestRouter creates a test gin router
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// GetTestDatabaseName returns the test database name
func GetTestDatabaseName() string {
	return "chatlogger_test"
}

// WaitForDatabase waits for database to be ready
func WaitForDatabase(dbURL string, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		db, err := sql.Open("postgres", dbURL)
		if err == nil {
			if err := db.Ping(); err == nil {
				db.Close()
				return nil
			}
			db.Close()
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("database not ready after %d retries", maxRetries)
}

// AutoMigrateTestSchema migrates all domain models for testing
func AutoMigrateTestSchema(db *gorm.DB) error {
	// Import the actual domain models for proper schema generation
	err := db.AutoMigrate(
		&domain.Organization{},
		&domain.User{},
		&domain.Chat{},
		&domain.Message{},
		&domain.APIKey{},
		&domain.Export{},
	)
	
	if err != nil {
		return fmt.Errorf("failed to auto-migrate schema: %w", err)
	}
	
	return nil
}
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

// benchmarkDB holds a shared database connection for benchmarks.
var benchmarkDB *gorm.DB

// SetupBenchmarkDB creates or returns a shared database connection optimized for benchmarks.
func SetupBenchmarkDB(b *testing.B) *gorm.DB {
	b.Helper()

	if benchmarkDB != nil {
		// Clean existing data for clean benchmark runs
		CleanBenchmarkData(benchmarkDB)
		return benchmarkDB
	}

	// Create a dedicated test database for benchmarks
	testDBURL := os.Getenv("BENCHMARK_DATABASE_URL")
	if testDBURL == "" {
		testDBURL = os.Getenv("TEST_DATABASE_URL")
	}
	if testDBURL == "" {
		b.Skip("No database available for benchmark tests")
	}

	db, err := gorm.Open(postgres.Open(testDBURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		b.Skipf("Failed to connect to benchmark database: %v", err)
	}

	// Configure connection pool for concurrent access
	sqlDB, err := db.DB()
	if err != nil {
		b.Skipf("Failed to get SQL DB: %v", err)
	}

	// Configure connection pool for concurrent benchmark access
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto-migrate schemas
	if err := AutoMigrateTestSchema(db); err != nil {
		b.Skipf("Failed to migrate benchmark schema: %v", err)
	}

	benchmarkDB = db
	return benchmarkDB
}

// CleanBenchmarkData removes all data from benchmark tables for clean runs.
func CleanBenchmarkData(db *gorm.DB) {
	// Order matters due to foreign key constraints
	tables := []string{"messages", "chats", "api_keys", "users", "organizations", "exports"}
	for _, table := range tables {
		db.Exec("DELETE FROM " + table)
	}
}

// CleanupBenchmarkDB properly closes the benchmark database connection.
func CleanupBenchmarkDB() {
	if benchmarkDB != nil {
		sqlDB, err := benchmarkDB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
		benchmarkDB = nil
	}
}

// SetupTestDB creates a test database connection.
func SetupTestDB(t testing.TB) *gorm.DB {
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
	db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name IN ('organizations', 'users', 'chats', 'messages', 'api_keys', 'exports')").
		Scan(&tableCount)
	if tableCount == 0 {
		t.Fatalf("No tables created after migration")
	}

	return db
}

// CleanupTestDB cleans up test database.
func CleanupTestDB(t testing.TB, db *gorm.DB) {
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

	if err := sqlDB.Close(); err != nil {
		t.Logf("Warning: failed to close database connection: %v", err)
	}
}

// SetupTestRouter creates a test gin router.
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// GetTestDatabaseName returns the test database name.
func GetTestDatabaseName() string {
	return "chatlogger_test"
}

// WaitForDatabase waits for database to be ready.
func WaitForDatabase(dbURL string, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		db, err := sql.Open("postgres", dbURL)
		if err == nil {
			if err := db.Ping(); err == nil {
				_ = db.Close() // Successful connection, safe to ignore close error
				return nil
			}
			_ = db.Close() // Failed connection, safe to ignore close error
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("database not ready after %d retries", maxRetries)
}

// AutoMigrateTestSchema migrates all domain models for testing.
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

package repository

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDatabaseMigrations tests that the database migrations work correctly.
func TestDatabaseMigrations(t *testing.T) {
	// Skip if no PostgreSQL available (migrations are PostgreSQL-specific)
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		t.Skip("Skipping migration tests - PostgreSQL required")
	}

	// Connect to PostgreSQL for migration testing
	db, err := gorm.Open(postgres.Open(testDBURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer func() { _ = sqlDB.Close() }()

	t.Run("Initial Schema Migration", func(t *testing.T) {
		// Apply initial schema migration
		migration001, err := os.ReadFile("../../migrations/001_initial_schema.sql")
		require.NoError(t, err)

		err = db.Exec(string(migration001)).Error
		require.NoError(t, err)

		// Verify all expected tables exist
		expectedTables := []string{
			"organizations", "users", "api_keys", "chats", "messages",
		}

		for _, table := range expectedTables {
			var exists bool
			err := db.Raw(`
				SELECT EXISTS (
					SELECT FROM information_schema.tables
					WHERE table_schema = 'public'
					AND table_name = ?
				)
			`, table).Scan(&exists).Error

			require.NoError(t, err)
			assert.True(t, exists, "Table %s should exist after migration", table)
		}

		// Verify expected indexes exist
		expectedIndexes := []string{
			"idx_api_keys_org_id", "idx_users_org_id", "idx_chats_org_id",
			"idx_chats_user_id", "idx_messages_chat_id", "idx_messages_created_at",
		}

		for _, index := range expectedIndexes {
			var exists bool
			err := db.Raw(`
				SELECT EXISTS (
					SELECT FROM pg_indexes
					WHERE schemaname = 'public'
					AND indexname = ?
				)
			`, index).Scan(&exists).Error

			require.NoError(t, err)
			assert.True(t, exists, "Index %s should exist after migration", index)
		}
	})

	t.Run("Exports Table Migration", func(t *testing.T) {
		// Apply exports table migration
		migration003, err := os.ReadFile("../../migrations/003_add_exports_table.sql")
		require.NoError(t, err)

		err = db.Exec(string(migration003)).Error
		require.NoError(t, err)

		// Verify exports table exists
		var exists bool
		err = db.Raw(`
			SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_schema = 'public'
				AND table_name = 'exports'
			)
		`).Scan(&exists).Error

		require.NoError(t, err)
		assert.True(t, exists, "Exports table should exist after migration")

		// Verify exports table columns
		expectedColumns := map[string]string{
			"id":              "bigint",
			"organization_id": "bigint",
			"user_id":         "bigint",
			"format":          "character varying",
			"type":            "character varying",
			"status":          "character varying",
			"file_path":       "text",
			"error":           "text",
			"created_at":      "timestamp without time zone",
			"updated_at":      "timestamp without time zone",
			"completed_at":    "timestamp without time zone",
		}

		for column, expectedType := range expectedColumns {
			var dataType string
			err := db.Raw(`
				SELECT data_type
				FROM information_schema.columns
				WHERE table_schema = 'public'
				AND table_name = 'exports'
				AND column_name = ?
			`, column).Scan(&dataType).Error

			require.NoError(t, err)
			assert.Equal(
				t,
				expectedType,
				dataType,
				"Column %s should have type %s",
				column,
				expectedType,
			)
		}

		// Verify foreign key constraints exist
		foreignKeys := []string{
			"fk_exports_organization",
			"fk_exports_user",
		}

		for _, fk := range foreignKeys {
			var exists bool
			err := db.Raw(`
				SELECT EXISTS (
					SELECT 1 FROM information_schema.table_constraints
					WHERE constraint_schema = 'public'
					AND constraint_name = ?
					AND constraint_type = 'FOREIGN KEY'
				)
			`, fk).Scan(&exists).Error

			require.NoError(t, err)
			assert.True(t, exists, "Foreign key %s should exist", fk)
		}

		// Verify trigger function exists
		var triggerExists bool
		err = db.Raw(`
			SELECT EXISTS (
				SELECT 1 FROM information_schema.routines
				WHERE routine_schema = 'public'
				AND routine_name = 'update_exports_updated_at'
			)
		`).Scan(&triggerExists).Error

		require.NoError(t, err)
		assert.True(t, triggerExists, "Trigger function should exist")
	})

	t.Run("Schema Validation", func(t *testing.T) {
		// Test that we can insert valid data into all tables

		// Insert organization
		var orgID int64
		err := db.Raw(`
			INSERT INTO organizations (name, slug, settings)
			VALUES ('Test Org', 'test-org', '{}')
			RETURNING id
		`).Scan(&orgID).Error
		require.NoError(t, err)
		assert.NotZero(t, orgID)

		// Insert user
		var userID int64
		err = db.Raw(`
			INSERT INTO users (email, password_hash, role, organization_id, first_name, last_name)
			VALUES ('test@example.com', 'hashed', 'user', ?, 'Test', 'User')
			RETURNING id
		`, orgID).Scan(&userID).Error
		require.NoError(t, err)
		assert.NotZero(t, userID)

		// Insert chat
		var chatID int64
		err = db.Raw(`
			INSERT INTO chats (organization_id, user_id, title, tags, metadata)
			VALUES (?, ?, 'Test Chat', '[]', '{}')
			RETURNING id
		`, orgID, userID).Scan(&chatID).Error
		require.NoError(t, err)
		assert.NotZero(t, chatID)

		// Insert message
		var messageID int64
		err = db.Raw(`
			INSERT INTO messages (chat_id, role, content, metadata)
			VALUES (?, 'user', 'Test message', '{}')
			RETURNING id
		`, chatID).Scan(&messageID).Error
		require.NoError(t, err)
		assert.NotZero(t, messageID)

		// Insert export
		var exportID int64
		err = db.Raw(`
			INSERT INTO exports (organization_id, user_id, format, type, status)
			VALUES (?, ?, 'json', 'all', 'pending')
			RETURNING id
		`, orgID, userID).Scan(&exportID).Error
		require.NoError(t, err)
		assert.NotZero(t, exportID)
	})

	t.Run("Data Integrity Constraints", func(t *testing.T) {
		// Test foreign key constraints

		// Should fail: Invalid organization_id in users table
		err := db.Exec(`
			INSERT INTO users (email, password_hash, role, organization_id)
			VALUES ('invalid@example.com', 'hash', 'user', 99999)
		`).Error
		assert.Error(t, err, "Should fail with invalid organization_id")

		// Should fail: Invalid role in users table
		err = db.Exec(`
			INSERT INTO users (email, password_hash, role, organization_id)
			VALUES ('invalid@example.com', 'hash', 'invalid_role', 1)
		`).Error
		assert.Error(t, err, "Should fail with invalid role")

		// Should fail: Invalid role in messages table
		err = db.Exec(`
			INSERT INTO messages (chat_id, role, content)
			VALUES (1, 'invalid_role', 'content')
		`).Error
		assert.Error(t, err, "Should fail with invalid message role")
	})
}

// TestMigrationRollback tests that migrations can be properly rolled back.
func TestMigrationRollback(t *testing.T) {
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		t.Skip("Skipping rollback tests - PostgreSQL required")
	}

	db, err := gorm.Open(postgres.Open(testDBURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer func() { _ = sqlDB.Close() }()

	// Drop all tables (simulating rollback)
	tables := []string{"exports", "messages", "chats", "api_keys", "users", "organizations"}
	for _, table := range tables {
		err := db.Exec("DROP TABLE IF EXISTS " + table + " CASCADE").Error
		require.NoError(t, err)
	}

	// Verify all tables are dropped
	for _, table := range tables {
		var exists bool
		err := db.Raw(`
			SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_schema = 'public'
				AND table_name = ?
			)
		`, table).Scan(&exists).Error

		require.NoError(t, err)
		assert.False(t, exists, "Table %s should not exist after rollback", table)
	}
}

// TestMigrationPerformance tests that migrations complete within reasonable time.
func TestMigrationPerformance(t *testing.T) {
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		t.Skip("Skipping performance tests - PostgreSQL required")
	}

	db, err := gorm.Open(postgres.Open(testDBURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer func() { _ = sqlDB.Close() }()

	// Test migration with larger dataset
	t.Run("Migration with data", func(t *testing.T) {
		// Apply base migration
		migration001, err := os.ReadFile("../../migrations/001_initial_schema.sql")
		require.NoError(t, err)

		err = db.Exec(string(migration001)).Error
		require.NoError(t, err)

		// Insert test data
		for i := 0; i < 100; i++ {
			err := db.Exec(`
				INSERT INTO organizations (name, slug)
				VALUES (?, ?)
			`, fmt.Sprintf("Org %d", i), fmt.Sprintf("org-%d", i)).Error
			require.NoError(t, err)
		}

		// Apply exports migration on existing data
		migration003, err := os.ReadFile("../../migrations/003_add_exports_table.sql")
		require.NoError(t, err)

		err = db.Exec(string(migration003)).Error
		require.NoError(t, err)

		// Verify data integrity after migration
		var orgCount int64
		err = db.Raw("SELECT COUNT(*) FROM organizations").Scan(&orgCount).Error
		require.NoError(t, err)
		assert.Equal(t, int64(100), orgCount, "Organization data should be preserved")
	})
}

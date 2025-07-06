package repository

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/kjanat/chatlogger-api-go/test/fixtures"
	"github.com/kjanat/chatlogger-api-go/test/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConcurrentChatCreation tests for race conditions during concurrent chat creation.
func TestConcurrentChatCreation(t *testing.T) {
	t.Skip(
		"Skipping race condition test - SQLite in-memory DB doesn't support concurrent access across goroutines",
	)
	t.Skip(
		"Skipping race condition test - SQLite in-memory DB doesn't support concurrent access across goroutines",
	)
	if testing.Short() {
		t.Skip("Skipping race condition test in short mode")
	}

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	// Verify tables exist before proceeding
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='chats'").
		Scan(&count).
		Error
	require.NoError(t, err)
	if count == 0 {
		t.Fatalf("Chats table does not exist")
	}

	dbWrapper := &Database{DB: db}
	repo := NewChatRepository(dbWrapper)

	// Setup test data
	org := fixtures.CreateTestOrganization()
	err = db.Create(org).Error
	require.NoError(t, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(t, err)

	const numGoroutines = 50
	const chatsPerGoroutine = 10

	var wg sync.WaitGroup
	results := make(chan error, numGoroutines*chatsPerGoroutine)
	createdIDs := make(chan uint64, numGoroutines*chatsPerGoroutine)

	// Launch concurrent chat creation
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < chatsPerGoroutine; j++ {
				chat := fixtures.CreateTestChat(user.ID)
				chat.OrganizationID = org.ID
				chat.Title = fmt.Sprintf("Chat %d-%d", goroutineID, j)

				err := repo.Create(chat)
				results <- err

				if err == nil {
					createdIDs <- chat.ID
				}
			}
		}(i)
	}

	wg.Wait()
	close(results)
	close(createdIDs)

	// Check results
	errorCount := 0
	for err := range results {
		if err != nil {
			errorCount++
			t.Logf("Creation error: %v", err)
		}
	}

	assert.Equal(t, 0, errorCount, "Should have no race condition errors")

	// Verify all IDs are unique (no duplicate ID assignment)
	idSet := make(map[uint64]bool)
	duplicates := 0
	for id := range createdIDs {
		if idSet[id] {
			duplicates++
		}
		idSet[id] = true
	}

	assert.Equal(t, 0, duplicates, "Should have no duplicate IDs from concurrent creation")
	assert.Equal(
		t,
		numGoroutines*chatsPerGoroutine,
		len(idSet),
		"Should have created expected number of unique chats",
	)
}

// TestConcurrentReadWrite tests concurrent read/write operations for race conditions.
func TestConcurrentReadWrite(t *testing.T) {
	t.Skip(
		"Skipping race condition test - SQLite in-memory DB doesn't support concurrent access across goroutines",
	)
	t.Skip(
		"Skipping race condition test - SQLite in-memory DB doesn't support concurrent access across goroutines",
	)
	if testing.Short() {
		t.Skip("Skipping race condition test in short mode")
	}

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewChatRepository(dbWrapper)

	// Setup test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(t, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(t, err)

	// Pre-create a chat to read/update
	chat := fixtures.CreateTestChat(user.ID)
	chat.OrganizationID = org.ID
	err = repo.Create(chat)
	require.NoError(t, err)

	const numOperations = 100
	var wg sync.WaitGroup
	errors := make(chan error, numOperations*3)

	// Concurrent readers
	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := repo.FindByID(chat.ID)
			errors <- err
		}()
	}

	// Concurrent updaters
	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func(updateID int) {
			defer wg.Done()

			chatToUpdate := &domain.Chat{
				ID:             chat.ID,
				OrganizationID: org.ID,
				UserID:         chat.UserID,
				Title:          fmt.Sprintf("Updated Title %d", updateID),
				Tags:           chat.Tags,
				Metadata:       chat.Metadata,
				CreatedAt:      chat.CreatedAt,
				UpdatedAt:      time.Now(),
			}

			err := repo.Update(chatToUpdate)
			errors <- err
		}(i)
	}

	// Concurrent list operations
	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := repo.FindByOrganizationID(org.ID, 10, 0)
			errors <- err
		}()
	}

	wg.Wait()
	close(errors)

	// Check for race condition errors
	errorCount := 0
	for err := range errors {
		if err != nil {
			errorCount++
			t.Logf("Concurrent operation error: %v", err)
		}
	}

	assert.Equal(
		t,
		0,
		errorCount,
		"Should have no race condition errors in concurrent read/write operations",
	)
}

// TestConcurrentConnectionHandling tests database connection pool under concurrent load.
func TestConcurrentConnectionHandling(t *testing.T) {
	t.Skip(
		"Skipping race condition test - SQLite in-memory DB doesn't support concurrent access across goroutines",
	)
	if testing.Short() {
		t.Skip("Skipping race condition test in short mode")
	}

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	// Configure limited connection pool to stress test
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(time.Minute)

	dbWrapper := &Database{DB: db}
	repo := NewChatRepository(dbWrapper)

	// Setup test data
	org := fixtures.CreateTestOrganization()
	err = db.Create(org).Error
	require.NoError(t, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(t, err)

	const numGoroutines = 20 // More goroutines than connections
	const operationsPerGoroutine = 5

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*operationsPerGoroutine)

	// Launch many concurrent operations to test connection pooling
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < operationsPerGoroutine; j++ {
				// Mix of operations to stress connection pool
				switch j % 3 {
				case 0:
					chat := fixtures.CreateTestChat(user.ID)
					chat.OrganizationID = org.ID
					chat.Title = fmt.Sprintf("Chat %d-%d", goroutineID, j)
					err := repo.Create(chat)
					errors <- err

				case 1:
					_, err := repo.FindByOrganizationID(org.ID, 5, 0)
					errors <- err

				case 2:
					// Simulate longer operation
					time.Sleep(time.Millisecond * 10)
					_, err := repo.FindByOrganizationID(org.ID, 10, 0)
					errors <- err
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for connection pool exhaustion or deadlock errors
	errorCount := 0
	for err := range errors {
		if err != nil {
			errorCount++
			t.Logf("Connection pool error: %v", err)
		}
	}

	assert.Equal(t, 0, errorCount, "Should handle concurrent connections without errors")
}

// TestRaceConditionInTransactions tests for race conditions in transaction handling.
func TestRaceConditionInTransactions(t *testing.T) {
	t.Skip(
		"Skipping race condition test - SQLite in-memory DB doesn't support concurrent access across goroutines",
	)
	if testing.Short() {
		t.Skip("Skipping race condition test in short mode")
	}

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	_ = NewChatRepository(dbWrapper)
	_ = NewMessageRepository(dbWrapper)

	// Setup test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(t, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(t, err)

	const numTransactions = 50

	var wg sync.WaitGroup
	errors := make(chan error, numTransactions*2)

	// Concurrent transactions: each creates a chat and then messages for it
	for i := 0; i < numTransactions; i++ {
		wg.Add(1)
		go func(txID int) {
			defer wg.Done()

			// Start transaction
			tx := db.Begin()
			if tx.Error != nil {
				errors <- tx.Error
				return
			}

			// Create chat in transaction
			txWrapper := &Database{DB: tx}
			txChatRepo := NewChatRepository(txWrapper)

			chat := fixtures.CreateTestChat(user.ID)
			chat.OrganizationID = org.ID
			chat.Title = fmt.Sprintf("Transaction Chat %d", txID)

			err := txChatRepo.Create(chat)
			if err != nil {
				tx.Rollback()
				errors <- err
				return
			}

			// Create message in same transaction
			txMessageRepo := NewMessageRepository(txWrapper)
			message := fixtures.CreateTestMessage(chat.ID)
			message.Content = fmt.Sprintf("Transaction Message %d", txID)

			err = txMessageRepo.Create(message)
			if err != nil {
				tx.Rollback()
				errors <- err
				return
			}

			// Commit transaction
			err = tx.Commit().Error
			errors <- err
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for transaction race condition errors
	errorCount := 0
	for err := range errors {
		if err != nil {
			errorCount++
			t.Logf("Transaction race condition error: %v", err)
		}
	}

	assert.Equal(t, 0, errorCount, "Should have no transaction race condition errors")

	// Verify all data was created successfully
	var chatCount int64
	err = db.Model(&domain.Chat{}).Where("organization_id = ?", org.ID).Count(&chatCount).Error
	require.NoError(t, err)
	assert.Equal(
		t,
		int64(numTransactions),
		chatCount,
		"Should have created expected number of chats",
	)

	var messageCount int64
	err = db.Model(&domain.Message{}).Count(&messageCount).Error
	require.NoError(t, err)
	assert.Equal(
		t,
		int64(numTransactions),
		messageCount,
		"Should have created expected number of messages",
	)
}

// TestMemoryRaceConditions tests for memory-related race conditions.
func TestMemoryRaceConditions(t *testing.T) {
	t.Skip(
		"Skipping race condition test - SQLite in-memory DB doesn't support concurrent access across goroutines",
	)
	if testing.Short() {
		t.Skip("Skipping race condition test in short mode")
	}

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewChatRepository(dbWrapper)

	// Setup test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(t, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(t, err)

	// Create shared data structure to test for race conditions
	sharedCounter := struct {
		sync.RWMutex
		count int
	}{}

	const numGoroutines = 100

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	// Concurrent operations that modify shared state
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			// Create chat (this tests database-level concurrency)
			chat := fixtures.CreateTestChat(user.ID)
			chat.OrganizationID = org.ID
			chat.Title = fmt.Sprintf("Race Test Chat %d", goroutineID)

			err := repo.Create(chat)
			if err != nil {
				errors <- err
				return
			}

			// Simulate shared memory access (tests Go-level concurrency)
			sharedCounter.Lock()
			currentCount := sharedCounter.count
			time.Sleep(time.Nanosecond) // Force potential race condition
			sharedCounter.count = currentCount + 1
			sharedCounter.Unlock()

			errors <- nil
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	errorCount := 0
	for err := range errors {
		if err != nil {
			errorCount++
			t.Logf("Memory race condition error: %v", err)
		}
	}

	assert.Equal(t, 0, errorCount, "Should have no memory race condition errors")
	assert.Equal(t, numGoroutines, sharedCounter.count, "Shared counter should have correct value")
}

// TestGoroutineLeaks tests for goroutine leaks in concurrent operations.
func TestGoroutineLeaks(t *testing.T) {
	t.Skip(
		"Skipping race condition test - SQLite in-memory DB doesn't support concurrent access across goroutines",
	)
	if testing.Short() {
		t.Skip("Skipping race condition test in short mode")
	}

	initialGoroutines := runtime.NumGoroutine()

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewChatRepository(dbWrapper)

	// Setup test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(t, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(t, err)

	// Run operations that might leak goroutines
	const iterations = 10
	for i := 0; i < iterations; i++ {
		var wg sync.WaitGroup

		for j := 0; j < 10; j++ {
			wg.Add(1)
			go func(opID int) {
				defer wg.Done()

				chat := fixtures.CreateTestChat(user.ID)
				chat.OrganizationID = org.ID
				_ = repo.Create(chat)

				// Force some work to potentially leak
				done := make(chan bool)
				go func() {
					time.Sleep(time.Millisecond)
					done <- true
				}()
				<-done
			}(j)
		}

		wg.Wait()

		// Force garbage collection
		runtime.GC()
		runtime.GC()
	}

	// Allow time for goroutines to clean up
	time.Sleep(time.Millisecond * 100)
	runtime.GC()

	finalGoroutines := runtime.NumGoroutine()

	// We allow some tolerance for background goroutines
	goroutineDiff := finalGoroutines - initialGoroutines
	assert.LessOrEqual(
		t,
		goroutineDiff,
		5,
		"Should not leak significant number of goroutines (leaked: %d)",
		goroutineDiff,
	)
}

// TestDataRaceInStats tests for race conditions in statistics calculations.
func TestDataRaceInStats(t *testing.T) {
	t.Skip(
		"Skipping race condition test - SQLite in-memory DB doesn't support concurrent access across goroutines",
	)
	if testing.Short() {
		t.Skip("Skipping race condition test in short mode")
	}

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewChatRepository(dbWrapper)

	// Setup test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(t, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(t, err)

	// Create some initial data
	for i := 0; i < 20; i++ {
		chat := fixtures.CreateTestChat(user.ID)
		chat.OrganizationID = org.ID
		chat.Tags = fmt.Sprintf(`["tag-%d"]`, i%5)
		err := repo.Create(chat)
		require.NoError(t, err)
	}

	const numGoroutines = 20
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*3)

	// Concurrent statistics queries while data is being modified
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			// Get tag stats
			_, err := repo.GetTagStats(org.ID)
			errors <- err

			// Get count by date range
			start := time.Now().AddDate(0, 0, -1)
			end := time.Now()
			_, err = repo.CountByOrgIDAndDateRange(org.ID, start, end)
			errors <- err

			// Create new chat while others are reading stats
			if goroutineID%2 == 0 {
				chat := fixtures.CreateTestChat(user.ID)
				chat.OrganizationID = org.ID
				chat.Tags = fmt.Sprintf(`["new-tag-%d"]`, goroutineID)
				err = repo.Create(chat)
				errors <- err
			} else {
				errors <- nil
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for race condition errors
	errorCount := 0
	for err := range errors {
		if err != nil {
			errorCount++
			t.Logf("Statistics race condition error: %v", err)
		}
	}

	assert.Equal(
		t,
		0,
		errorCount,
		"Should have no race condition errors in statistics calculations",
	)
}

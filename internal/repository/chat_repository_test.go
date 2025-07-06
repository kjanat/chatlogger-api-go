package repository

import (
	"testing"
	"time"

	"github.com/kjanat/chatlogger-api-go/test/fixtures"
	"github.com/kjanat/chatlogger-api-go/test/testutils"
	"github.com/stretchr/testify/assert"
)

func TestChatRepository_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewChatRepository(dbWrapper)

	// Create test organization first
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	testutils.SetupError(t, err, "creating test organization")

	// Create test user
	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	testutils.SetupError(t, err, "creating test user")

	// Create test chat
	originalTime := time.Now()
	chat := fixtures.CreateTestChat(user.ID)
	chat.OrganizationID = org.ID

	err = repo.Create(chat)
	testutils.RepositoryError(t, err, "ChatRepository", "Create")

	// Comprehensive validation
	assert.NotZero(t, chat.ID, "Chat ID should be auto-generated")
	assert.Equal(t, org.ID, chat.OrganizationID, "Organization ID should match")
	assert.Equal(t, user.ID, *chat.UserID, "User ID should match")
	assert.NotEmpty(t, chat.Title, "Chat title should not be empty")
	assert.NotEmpty(t, chat.Tags, "Chat tags should be set")
	assert.NotEmpty(t, chat.Metadata, "Chat metadata should be set")

	// Validate timestamps
	assert.True(
		t,
		chat.CreatedAt.After(originalTime.Add(-1*time.Second)),
		"CreatedAt should be recent",
	)
	assert.True(
		t,
		chat.UpdatedAt.After(originalTime.Add(-1*time.Second)),
		"UpdatedAt should be recent",
	)
	assert.WithinDuration(
		t,
		chat.CreatedAt,
		chat.UpdatedAt,
		time.Millisecond*100,
		"CreatedAt and UpdatedAt should be close",
	)

	// Validate chat was actually persisted
	var count int64
	err = db.Model(chat).Where("id = ?", chat.ID).Count(&count).Error
	testutils.DatabaseError(t, err, "count validation", "chats")
	assert.Equal(t, int64(1), count, "Chat should be persisted in database")
}

func TestChatRepository_FindByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewChatRepository(dbWrapper)

	// Create test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	testutils.SetupError(t, err, "creating test organization")

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	testutils.SetupError(t, err, "creating test user")

	chat := fixtures.CreateTestChat(user.ID)
	chat.OrganizationID = org.ID
	err = repo.Create(chat)
	testutils.SetupError(t, err, "creating test chat for FindByID test")

	// Test finding existing chat
	foundChat, err := repo.FindByID(chat.ID)
	testutils.RepositoryError(t, err, "ChatRepository", "FindByID")
	assert.NotNil(t, foundChat, "Found chat should not be nil")

	// Comprehensive validation of retrieved chat
	assert.Equal(t, chat.ID, foundChat.ID, "Chat ID should match")
	assert.Equal(t, chat.Title, foundChat.Title, "Chat title should match")
	assert.Equal(t, chat.OrganizationID, foundChat.OrganizationID, "Organization ID should match")
	assert.Equal(t, chat.UserID, foundChat.UserID, "User ID should match")
	assert.Equal(t, chat.Tags, foundChat.Tags, "Tags should match")
	assert.Equal(t, chat.Metadata, foundChat.Metadata, "Metadata should match")
	assert.WithinDuration(
		t,
		chat.CreatedAt,
		foundChat.CreatedAt,
		time.Second,
		"CreatedAt should match",
	)
	assert.WithinDuration(
		t,
		chat.UpdatedAt,
		foundChat.UpdatedAt,
		time.Second,
		"UpdatedAt should match",
	)

	// Test finding non-existent chat
	nonExistentChat, err := repo.FindByID(9999)
	testutils.RepositoryError(t, err, "ChatRepository", "FindByID for non-existent chat")
	assert.Nil(t, nonExistentChat, "Non-existent chat should return nil")

	// Test edge cases
	zeroIDChat, err := repo.FindByID(0)
	testutils.RepositoryError(t, err, "ChatRepository", "FindByID with zero ID")
	assert.Nil(t, zeroIDChat, "Zero ID should return nil")
}

func TestChatRepository_FindByOrganizationID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewChatRepository(dbWrapper)

	// Create test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	testutils.SetupError(t, err, "creating test organization")

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	testutils.SetupError(t, err, "creating test user")

	// Create multiple chats
	chat1 := fixtures.CreateTestChat(user.ID)
	chat1.OrganizationID = org.ID
	chat1.Title = "Chat 1"
	err = repo.Create(chat1)
	testutils.SetupError(t, err, "creating first test chat")

	chat2 := fixtures.CreateTestChat(user.ID)
	chat2.OrganizationID = org.ID
	chat2.Title = "Chat 2"
	err = repo.Create(chat2)
	testutils.SetupError(t, err, "creating second test chat")

	// Test finding chats by organization ID
	chats, err := repo.FindByOrganizationID(org.ID, 10, 0)
	testutils.RepositoryError(t, err, "ChatRepository", "FindByOrganizationID")
	assert.Len(t, chats, 2)

	// Test pagination
	chats, err = repo.FindByOrganizationID(org.ID, 1, 0)
	testutils.RepositoryError(t, err, "ChatRepository", "FindByOrganizationID with limit 1")
	assert.Len(t, chats, 1)

	chats, err = repo.FindByOrganizationID(org.ID, 1, 1)
	testutils.RepositoryError(t, err, "ChatRepository", "FindByOrganizationID with offset 1")
	assert.Len(t, chats, 1)
}

func TestChatRepository_FindByUserID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewChatRepository(dbWrapper)

	// Create test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	testutils.SetupError(t, err, "creating test organization")

	user1 := fixtures.CreateTestUser(org.ID)
	user1.Email = "user1@example.com"
	err = db.Create(user1).Error
	testutils.SetupError(t, err, "creating test user1")

	user2 := fixtures.CreateTestUser(org.ID)
	user2.Email = "user2@example.com"
	err = db.Create(user2).Error
	testutils.SetupError(t, err, "creating test user2")

	// Create chats for user1
	chat1 := fixtures.CreateTestChat(user1.ID)
	chat1.OrganizationID = org.ID
	err = repo.Create(chat1)
	testutils.SetupError(t, err, "creating chat1 for user1")

	chat2 := fixtures.CreateTestChat(user1.ID)
	chat2.OrganizationID = org.ID
	err = repo.Create(chat2)
	testutils.SetupError(t, err, "creating chat2 for user1")

	// Create chat for user2
	chat3 := fixtures.CreateTestChat(user2.ID)
	chat3.OrganizationID = org.ID
	err = repo.Create(chat3)
	testutils.SetupError(t, err, "creating chat for user2")

	// Test finding chats by user ID
	user1Chats, err := repo.FindByUserID(user1.ID, 10, 0)
	testutils.RepositoryError(t, err, "ChatRepository", "FindByUserID for user1")
	assert.Len(t, user1Chats, 2)

	user2Chats, err := repo.FindByUserID(user2.ID, 10, 0)
	testutils.RepositoryError(t, err, "ChatRepository", "FindByUserID for user2")
	assert.Len(t, user2Chats, 1)
}

func TestChatRepository_Update(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewChatRepository(dbWrapper)

	// Create test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	testutils.SetupError(t, err, "creating test organization")

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	testutils.SetupError(t, err, "creating test user")

	chat := fixtures.CreateTestChat(user.ID)
	chat.OrganizationID = org.ID
	err = repo.Create(chat)
	testutils.SetupError(t, err, "creating test chat for update test")

	// Update chat
	chat.Title = "Updated Title"
	chat.Title = "Updated Title"
	err = repo.Update(chat)
	testutils.RepositoryError(t, err, "ChatRepository", "Update")

	// Verify update
	updatedChat, err := repo.FindByID(chat.ID)
	testutils.RepositoryError(t, err, "ChatRepository", "FindByID after update")
	assert.Equal(t, "Updated Title", updatedChat.Title)
	assert.Equal(t, "Updated Title", updatedChat.Title)
}

func TestChatRepository_Delete(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewChatRepository(dbWrapper)

	// Create test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	testutils.SetupError(t, err, "creating test organization")

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	testutils.SetupError(t, err, "creating test user")

	chat := fixtures.CreateTestChat(user.ID)
	chat.OrganizationID = org.ID
	err = repo.Create(chat)
	testutils.SetupError(t, err, "creating test chat for delete test")

	// Delete chat
	err = repo.Delete(chat.ID)
	testutils.RepositoryError(t, err, "ChatRepository", "Delete")

	// Verify deletion
	deletedChat, err := repo.FindByID(chat.ID)
	testutils.RepositoryError(t, err, "ChatRepository", "FindByID after delete")
	assert.Nil(t, deletedChat)
}

func TestChatRepository_CountByOrgIDAndDateRange(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewChatRepository(dbWrapper)

	// Create test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	testutils.SetupError(t, err, "creating test organization")

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	testutils.SetupError(t, err, "creating test user")

	// Create chats with different timestamps
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)

	chat1 := fixtures.CreateTestChat(user.ID)
	chat1.OrganizationID = org.ID
	chat1.CreatedAt = yesterday
	err = repo.Create(chat1)
	testutils.SetupError(t, err, "creating chat1 with yesterday timestamp")

	chat2 := fixtures.CreateTestChat(user.ID)
	chat2.OrganizationID = org.ID
	chat2.CreatedAt = now
	err = repo.Create(chat2)
	testutils.SetupError(t, err, "creating chat2 with current timestamp")

	// Count chats in range
	count, err := repo.CountByOrgIDAndDateRange(org.ID, yesterday.Add(-time.Hour), tomorrow)
	testutils.RepositoryError(t, err, "ChatRepository", "CountByOrgIDAndDateRange wide range")
	assert.Equal(t, int64(2), count)

	// Count chats in narrow range
	count, err = repo.CountByOrgIDAndDateRange(org.ID, now.Add(-time.Hour), now.Add(time.Hour))
	testutils.RepositoryError(t, err, "ChatRepository", "CountByOrgIDAndDateRange narrow range")
	assert.Equal(t, int64(1), count)
}

func TestChatRepository_GetTagStats(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewChatRepository(dbWrapper)

	// Create test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	testutils.SetupError(t, err, "creating test organization")

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	testutils.SetupError(t, err, "creating test user")

	chat := fixtures.CreateTestChat(user.ID)
	chat.OrganizationID = org.ID
	err = repo.Create(chat)
	testutils.SetupError(t, err, "creating test chat for tag stats test")

	// Test getting tag stats (simplified implementation)
	stats, err := repo.GetTagStats(org.ID)
	testutils.RepositoryError(t, err, "ChatRepository", "GetTagStats")
	assert.NotNil(t, stats)
}

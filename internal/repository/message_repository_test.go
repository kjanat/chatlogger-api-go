package repository

import (
	"testing"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/kjanat/chatlogger-api-go/test/fixtures"
	"github.com/kjanat/chatlogger-api-go/test/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageRepository_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewMessageRepository(dbWrapper)

	// Create test data hierarchy
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(t, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(t, err)

	chat := fixtures.CreateTestChat(user.ID)
	chat.OrganizationID = org.ID
	err = db.Create(chat).Error
	require.NoError(t, err)

	// Create test message
	message := fixtures.CreateTestMessage(chat.ID)

	err = repo.Create(message)
	assert.NoError(t, err)
	assert.NotZero(t, message.ID)
	assert.Equal(t, chat.ID, message.ChatID)
}

func TestMessageRepository_FindByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewMessageRepository(dbWrapper)

	// Create test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(t, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(t, err)

	chat := fixtures.CreateTestChat(user.ID)
	chat.OrganizationID = org.ID
	err = db.Create(chat).Error
	require.NoError(t, err)

	message := fixtures.CreateTestMessage(chat.ID)
	err = repo.Create(message)
	require.NoError(t, err)

	// Test finding existing message
	foundMessage, err := repo.FindByID(message.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundMessage)
	assert.Equal(t, message.ID, foundMessage.ID)
	assert.Equal(t, message.Content, foundMessage.Content)
	assert.Equal(t, message.Role, foundMessage.Role)

	// Test finding non-existent message
	nonExistentMessage, err := repo.FindByID(9999)
	assert.NoError(t, err)
	assert.Nil(t, nonExistentMessage)
}

func TestMessageRepository_FindByChatID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewMessageRepository(dbWrapper)

	// Create test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(t, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(t, err)

	chat1 := fixtures.CreateTestChat(user.ID)
	chat1.OrganizationID = org.ID
	chat1.Title = "Chat 1"
	err = db.Create(chat1).Error
	require.NoError(t, err)

	chat2 := fixtures.CreateTestChat(user.ID)
	chat2.OrganizationID = org.ID
	chat2.Title = "Chat 2"
	err = db.Create(chat2).Error
	require.NoError(t, err)

	// Create messages for chat1
	message1 := fixtures.CreateTestMessage(chat1.ID)
	message1.Content = "First message"
	message1.Role = domain.MessageRoleUser
	message1.CreatedAt = time.Now().Add(-2 * time.Hour)
	err = repo.Create(message1)
	require.NoError(t, err)

	message2 := fixtures.CreateTestMessage(chat1.ID)
	message2.Content = "Second message"
	message2.Role = domain.MessageRoleAssistant
	message2.CreatedAt = time.Now().Add(-1 * time.Hour)
	err = repo.Create(message2)
	require.NoError(t, err)

	// Create message for chat2
	message3 := fixtures.CreateTestMessage(chat2.ID)
	message3.Content = "Chat 2 message"
	err = repo.Create(message3)
	require.NoError(t, err)

	// Test finding messages by chat ID
	chat1Messages, err := repo.FindByChatID(chat1.ID)
	assert.NoError(t, err)
	assert.Len(t, chat1Messages, 2)

	// Verify ordering (should be ASC by created_at)
	assert.Equal(t, "First message", chat1Messages[0].Content)
	assert.Equal(t, "Second message", chat1Messages[1].Content)

	chat2Messages, err := repo.FindByChatID(chat2.ID)
	assert.NoError(t, err)
	assert.Len(t, chat2Messages, 1)
	assert.Equal(t, "Chat 2 message", chat2Messages[0].Content)
}

func TestMessageRepository_CountByOrgIDAndDateRange(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewMessageRepository(dbWrapper)

	// Create test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(t, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(t, err)

	chat := fixtures.CreateTestChat(user.ID)
	chat.OrganizationID = org.ID
	err = db.Create(chat).Error
	require.NoError(t, err)

	// Create messages with different timestamps
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)

	message1 := fixtures.CreateTestMessage(chat.ID)
	message1.CreatedAt = yesterday
	err = repo.Create(message1)
	require.NoError(t, err)

	message2 := fixtures.CreateTestMessage(chat.ID)
	message2.CreatedAt = now
	err = repo.Create(message2)
	require.NoError(t, err)

	message3 := fixtures.CreateTestMessage(chat.ID)
	message3.CreatedAt = now.Add(1 * time.Hour)
	err = repo.Create(message3)
	require.NoError(t, err)

	// Count messages in range
	count, err := repo.CountByOrgIDAndDateRange(org.ID, yesterday.Add(-time.Hour), tomorrow)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)

	// Count messages in narrow range
	count, err = repo.CountByOrgIDAndDateRange(org.ID, now.Add(-time.Hour), now.Add(30*time.Minute))
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestMessageRepository_GetRoleStats(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewMessageRepository(dbWrapper)

	// Create test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(t, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(t, err)

	chat := fixtures.CreateTestChat(user.ID)
	chat.OrganizationID = org.ID
	err = db.Create(chat).Error
	require.NoError(t, err)

	// Create messages with different roles
	userMessage1 := fixtures.CreateTestMessage(chat.ID)
	userMessage1.Role = domain.MessageRoleUser
	err = repo.Create(userMessage1)
	require.NoError(t, err)

	userMessage2 := fixtures.CreateTestMessage(chat.ID)
	userMessage2.Role = domain.MessageRoleUser
	err = repo.Create(userMessage2)
	require.NoError(t, err)

	assistantMessage := fixtures.CreateTestMessage(chat.ID)
	assistantMessage.Role = domain.MessageRoleAssistant
	err = repo.Create(assistantMessage)
	require.NoError(t, err)

	systemMessage := fixtures.CreateTestMessage(chat.ID)
	systemMessage.Role = domain.MessageRoleSystem
	err = repo.Create(systemMessage)
	require.NoError(t, err)

	// Get role statistics
	stats, err := repo.GetRoleStats(org.ID)
	assert.NoError(t, err)
	assert.NotNil(t, stats)

	expectedStats := map[domain.MessageRole]int64{
		domain.MessageRoleUser:      2,
		domain.MessageRoleAssistant: 1,
		domain.MessageRoleSystem:    1,
	}

	assert.Equal(t, expectedStats, stats)
}

func TestMessageRepository_MessageValidation(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewMessageRepository(dbWrapper)

	// Create test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(t, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(t, err)

	chat := fixtures.CreateTestChat(user.ID)
	chat.OrganizationID = org.ID
	err = db.Create(chat).Error
	require.NoError(t, err)

	// Test creating message with all valid roles
	validRoles := []domain.MessageRole{
		domain.MessageRoleUser,
		domain.MessageRoleAssistant,
		domain.MessageRoleSystem,
	}

	for _, role := range validRoles {
		message := fixtures.CreateTestMessage(chat.ID)
		message.Role = role
		message.Content = "Test content for " + string(role)

		err = repo.Create(message)
		assert.NoError(t, err, "Should be able to create message with role: %s", role)
	}
}

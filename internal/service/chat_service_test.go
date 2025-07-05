package service

import (
	"errors"
	"testing"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/kjanat/chatlogger-api-go/test/fixtures"
	"github.com/kjanat/chatlogger-api-go/test/mocks"
	"github.com/kjanat/chatlogger-api-go/test/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChatService_CreateChat(t *testing.T) {
	// Using new mock builder pattern
	mockRepo := mocks.NewMockChatRepositoryBuilder().
		WithCreateReturns(nil).
		Build()
	service := NewChatService(mockRepo)

	chat := &domain.Chat{
		OrganizationID: 1,
		UserID:         uintPtr(1),
		Title:          "Test Chat",
		Tags:           `["test"]`,
		Metadata:       `{"test": true}`,
	}

	err := service.CreateChat(chat)

	testutils.ServiceError(t, err, "ChatService", "CreateChat")
	assert.NotZero(t, chat.CreatedAt)
	assert.NotZero(t, chat.UpdatedAt)
	mockRepo.AssertExpectations(t)
}

func TestChatService_CreateChat_RepositoryError(t *testing.T) {
	chat := fixtures.CreateTestChat(1)
	expectedError := errors.New("database error")

	// Using builder pattern for error scenario
	mockRepo := mocks.NewMockChatRepositoryBuilder().
		WithCreateReturns(expectedError).
		Build()
	service := NewChatService(mockRepo)

	err := service.CreateChat(chat)

	testutils.ExpectErrorWithMessage(t, err, "database error", "CreateChat with repository error")
	mockRepo.AssertExpectations(t)
}

func TestChatService_GetByID(t *testing.T) {
	expectedChat := fixtures.CreateTestChat(1)
	expectedChat.ID = 1

	// Using builder pattern with specific return values
	mockRepo := mocks.NewMockChatRepositoryBuilder().
		WithFindByIDReturns(expectedChat, nil).
		Build()
	service := NewChatService(mockRepo)

	chat, err := service.GetByID(1)

	testutils.ServiceError(t, err, "ChatService", "GetByID")
	assert.Equal(t, expectedChat, chat)
	mockRepo.AssertExpectations(t)
}

func TestChatService_GetByID_NotFound(t *testing.T) {
	// Using scenario builder for common case
	scenario := mocks.NewScenarioBuilder()
	mockRepo := scenario.ChatNotFoundScenario()
	service := NewChatService(mockRepo)

	chat, err := service.GetByID(999)

	assert.NoError(t, err)
	assert.Nil(t, chat)
	mockRepo.AssertExpectations(t)
}

func TestChatService_GetByOrganizationID(t *testing.T) {
	expectedChats := []domain.Chat{
		*fixtures.CreateTestChat(1),
		*fixtures.CreateTestChat(1),
	}

	// Using builder pattern for list operations
	mockRepo := mocks.NewMockChatRepositoryBuilder().
		WithFindByOrganizationIDReturns(expectedChats, nil).
		Build()
	service := NewChatService(mockRepo)

	chats, err := service.GetByOrganizationID(1, 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, expectedChats, chats)
	mockRepo.AssertExpectations(t)
}

func TestChatService_GetByUserID(t *testing.T) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	expectedChats := []domain.Chat{
		*fixtures.CreateTestChat(1),
	}

	mockRepo.On("FindByUserID", uint64(1), 5, 0).Return(expectedChats, nil)

	chats, err := service.GetByUserID(1, 5, 0)

	assert.NoError(t, err)
	assert.Equal(t, expectedChats, chats)
	mockRepo.AssertExpectations(t)
}

func TestChatService_UpdateChat(t *testing.T) {
	existingChat := fixtures.CreateTestChat(1)
	existingChat.ID = 1

	chatToUpdate := &domain.Chat{
		ID:             1,
		OrganizationID: 1,
		UserID:         uintPtr(1),
		Title:          "Updated Title",
		Tags:           `["updated"]`,
		Metadata:       `{"updated": true}`,
	}

	// Using builder pattern to chain multiple mock setups
	mockRepo := mocks.NewMockChatRepositoryBuilder().
		WithFindByIDReturns(existingChat, nil).
		WithUpdateReturns(nil).
		Build()
	service := NewChatService(mockRepo)

	err := service.UpdateChat(chatToUpdate)

	assert.NoError(t, err)
	assert.NotZero(t, chatToUpdate.UpdatedAt)
	mockRepo.AssertExpectations(t)
}

func TestChatService_UpdateChat_NotFound(t *testing.T) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	chatToUpdate := fixtures.CreateTestChat(1)
	chatToUpdate.ID = 999

	mockRepo.On("FindByID", uint64(999)).Return((*domain.Chat)(nil), nil)

	err := service.UpdateChat(chatToUpdate)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "chat not found")
	mockRepo.AssertExpectations(t)
}

func TestChatService_UpdateChat_FindError(t *testing.T) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	chatToUpdate := fixtures.CreateTestChat(1)
	chatToUpdate.ID = 1
	expectedError := errors.New("database error")

	mockRepo.On("FindByID", uint64(1)).Return((*domain.Chat)(nil), expectedError)

	err := service.UpdateChat(chatToUpdate)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error finding chat")
	mockRepo.AssertExpectations(t)
}

func TestChatService_DeleteChat(t *testing.T) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	mockRepo.On("Delete", uint64(1)).Return(nil)

	err := service.DeleteChat(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestChatService_DeleteChat_Error(t *testing.T) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	expectedError := errors.New("delete error")
	mockRepo.On("Delete", uint64(1)).Return(expectedError)

	err := service.DeleteChat(1)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestChatService_GetChatStats(t *testing.T) {
	orgID := uint64(1)
	start := time.Now().AddDate(0, 0, -7)
	end := time.Now()

	expectedChatCount := int64(10)
	expectedTagStats := map[string]int64{
		"support":   5,
		"technical": 3,
		"general":   2,
	}

	// Using builder pattern for complex stats operations
	mockRepo := mocks.NewMockChatRepositoryBuilder().
		WithCountByOrgIDAndDateRangeReturns(expectedChatCount, nil).
		WithGetTagStatsReturns(expectedTagStats, nil).
		Build()
	service := NewChatService(mockRepo)

	stats, err := service.GetChatStats(orgID, start, end)

	testutils.ServiceError(t, err, "ChatService", "GetChatStats")
	assert.Equal(t, expectedChatCount, stats["total_chats"])
	assert.Equal(t, expectedTagStats, stats["tag_stats"])

	dateRange, ok := stats["date_range"].(map[string]string)
	assert.True(t, ok, "date_range should be a map[string]string")
	assert.NotEmpty(t, dateRange["start"])
	assert.NotEmpty(t, dateRange["end"])

	mockRepo.AssertExpectations(t)
}

func TestChatService_GetChatStats_CountError(t *testing.T) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	orgID := uint64(1)
	start := time.Now().AddDate(0, 0, -7)
	end := time.Now()
	expectedError := errors.New("count error")

	mockRepo.On("CountByOrgIDAndDateRange", orgID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(int64(0), expectedError)

	stats, err := service.GetChatStats(orgID, start, end)

	assert.Error(t, err)
	assert.Nil(t, stats)
	assert.Contains(t, err.Error(), "error getting chat count")
	mockRepo.AssertExpectations(t)
}

func TestChatService_GetChatStats_TagStatsError(t *testing.T) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	orgID := uint64(1)
	start := time.Now().AddDate(0, 0, -7)
	end := time.Now()
	expectedError := errors.New("tag stats error")

	mockRepo.On("CountByOrgIDAndDateRange", orgID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(int64(10), nil)
	mockRepo.On("GetTagStats", orgID).Return(map[string]int64(nil), expectedError)

	stats, err := service.GetChatStats(orgID, start, end)

	assert.Error(t, err)
	assert.Nil(t, stats)
	assert.Contains(t, err.Error(), "error getting tag stats")
	mockRepo.AssertExpectations(t)
}

// Example test showing integration with builder pattern for specific operations
func TestChatService_Integration_BuilderPattern(t *testing.T) {
	// Using the test setup with focused mock configuration
	setup := mocks.NewServiceTestSetup()
	mockRepo := mocks.NewMockChatRepositoryBuilder().
		WithCreateReturns(nil).
		WithFindByIDReturns(setup.Chat, nil).
		WithFindByOrganizationIDReturns([]domain.Chat{*setup.Chat}, nil).
		Build()
	service := NewChatService(mockRepo)

	// Test create operation
	err := service.CreateChat(setup.Chat)
	assert.NoError(t, err)

	// Test get operation
	chat, err := service.GetByID(setup.Chat.ID)
	assert.NoError(t, err)
	assert.NotNil(t, chat)
	assert.Equal(t, setup.Chat.ID, chat.ID)
	assert.Equal(t, setup.Chat.Title, chat.Title)

	// Test list operation
	chats, err := service.GetByOrganizationID(setup.Org.ID, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, chats, 1)
	assert.Equal(t, setup.Chat.ID, chats[0].ID)

	mockRepo.AssertExpectations(t)
}

// Helper function to create uint64 pointer
func uintPtr(i uint64) *uint64 {
	return &i
}

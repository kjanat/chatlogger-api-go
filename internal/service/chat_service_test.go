package service

import (
	"errors"
	"testing"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/kjanat/chatlogger-api-go/test/fixtures"
	"github.com/kjanat/chatlogger-api-go/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestChatService_CreateChat(t *testing.T) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	chat := &domain.Chat{
		OrganizationID: 1,
		UserID:         uintPtr(1),
		Title:          "Test Chat",
		Tags:           `["test"]`,
		Metadata:       `{"test": true}`,
	}

	mockRepo.On("Create", mock.AnythingOfType("*domain.Chat")).Return(nil)

	err := service.CreateChat(chat)

	assert.NoError(t, err)
	assert.NotZero(t, chat.CreatedAt)
	assert.NotZero(t, chat.UpdatedAt)
	mockRepo.AssertExpectations(t)
}

func TestChatService_CreateChat_RepositoryError(t *testing.T) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	chat := fixtures.CreateTestChat(1)
	expectedError := errors.New("database error")

	mockRepo.On("Create", mock.AnythingOfType("*domain.Chat")).Return(expectedError)

	err := service.CreateChat(chat)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestChatService_GetByID(t *testing.T) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	expectedChat := fixtures.CreateTestChat(1)
	expectedChat.ID = 1

	mockRepo.On("FindByID", uint64(1)).Return(expectedChat, nil)

	chat, err := service.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedChat, chat)
	mockRepo.AssertExpectations(t)
}

func TestChatService_GetByID_NotFound(t *testing.T) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	mockRepo.On("FindByID", uint64(999)).Return((*domain.Chat)(nil), nil)

	chat, err := service.GetByID(999)

	assert.NoError(t, err)
	assert.Nil(t, chat)
	mockRepo.AssertExpectations(t)
}

func TestChatService_GetByOrganizationID(t *testing.T) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	expectedChats := []domain.Chat{
		*fixtures.CreateTestChat(1),
		*fixtures.CreateTestChat(1),
	}

	mockRepo.On("FindByOrganizationID", uint64(1), 10, 0).Return(expectedChats, nil)

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
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

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

	mockRepo.On("FindByID", uint64(1)).Return(existingChat, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Chat")).Return(nil)

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
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	orgID := uint64(1)
	start := time.Now().AddDate(0, 0, -7)
	end := time.Now()
	
	expectedChatCount := int64(10)
	expectedTagStats := map[string]int64{
		"support":   5,
		"technical": 3,
		"general":   2,
	}

	mockRepo.On("CountByOrgIDAndDateRange", orgID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(expectedChatCount, nil)
	mockRepo.On("GetTagStats", orgID).Return(expectedTagStats, nil)

	stats, err := service.GetChatStats(orgID, start, end)

	require.NoError(t, err)
	assert.Equal(t, expectedChatCount, stats["total_chats"])
	assert.Equal(t, expectedTagStats, stats["tag_stats"])
	
	dateRange, ok := stats["date_range"].(map[string]string)
	require.True(t, ok)
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

// Helper function to create uint64 pointer
func uintPtr(i uint64) *uint64 {
	return &i
}
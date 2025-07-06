package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/kjanat/chatlogger-api-go/test/fixtures"
	"github.com/kjanat/chatlogger-api-go/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupChatHandler() (*ChatHandler, *mocks.MockChatService, *mocks.MockMessageService) {
	mockChatService := &mocks.MockChatService{}
	mockMessageService := &mocks.MockMessageService{}
	handler := NewChatHandler(mockChatService, mockMessageService)
	return handler, mockChatService, mockMessageService
}

func setupGinContext(
	method, path string,
	body interface{},
) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	return c, w
}

func TestChatHandler_CreateChat_Success(t *testing.T) {
	handler, mockChatService, _ := setupChatHandler()

	reqBody := CreateChatRequest{
		Title: "Test Chat",
		Tags:  []string{"test", "api"},
		Metadata: &domain.ChatMetadata{
			IPAddress:  "192.168.1.1",
			Sentiment:  "positive",
			TokenCount: 150,
		},
		UserID: uintPtr(1),
	}

	c, w := setupGinContext("POST", "/v1/chats", reqBody)
	c.Set("orgID", uint64(100))

	mockChatService.On("CreateChat", mock.AnythingOfType("*domain.Chat")).
		Return(nil).
		Run(func(args mock.Arguments) {
			chat := args.Get(0).(*domain.Chat)
			chat.ID = 1 // Simulate DB assignment
		})

	handler.CreateChat(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Chat created successfully", response["message"])
	assert.Equal(t, float64(1), response["chat_id"])

	mockChatService.AssertExpectations(t)
}

func TestChatHandler_CreateChat_InvalidJSON(t *testing.T) {
	handler, _, _ := setupChatHandler()

	c, w := setupGinContext("POST", "/v1/chats", nil)
	c.Request.Body = http.NoBody
	c.Set("orgID", uint64(100))

	handler.CreateChat(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["error"], "Invalid request data")
}

func TestChatHandler_CreateChat_NoOrgID(t *testing.T) {
	handler, _, _ := setupChatHandler()

	reqBody := CreateChatRequest{Title: "Test Chat"}
	c, w := setupGinContext("POST", "/v1/chats", reqBody)
	// Not setting orgID in context

	handler.CreateChat(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Organization ID not found in context", response["error"])
}

func TestChatHandler_CreateChat_ServiceError(t *testing.T) {
	handler, mockChatService, _ := setupChatHandler()

	reqBody := CreateChatRequest{Title: "Test Chat"}
	c, w := setupGinContext("POST", "/v1/chats", reqBody)
	c.Set("orgID", uint64(100))

	mockChatService.On("CreateChat", mock.AnythingOfType("*domain.Chat")).Return(assert.AnError)

	handler.CreateChat(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["error"], "Failed to create chat")

	mockChatService.AssertExpectations(t)
}

func TestChatHandler_GetChat_Success(t *testing.T) {
	handler, mockChatService, mockMessageService := setupChatHandler()

	expectedChat := fixtures.CreateTestChat(1)
	expectedChat.ID = 1
	expectedChat.OrganizationID = 100

	c, w := setupGinContext("GET", "/v1/chats/1", nil)
	c.Set("orgID", uint64(100))
	c.Params = []gin.Param{{Key: "chatID", Value: "1"}}

	mockChatService.On("GetByID", uint64(1)).Return(expectedChat, nil)

	handler.GetChat(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response GetChatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedChat.ID, response.ID)
	assert.Equal(t, expectedChat.Title, response.Title)
	assert.NotNil(t, response.ParsedTags)
	assert.NotNil(t, response.ParsedMetadata)

	mockChatService.AssertExpectations(t)
	mockMessageService.AssertExpectations(t)
}

func TestChatHandler_GetChat_NotFound(t *testing.T) {
	handler, mockChatService, _ := setupChatHandler()

	c, w := setupGinContext("GET", "/v1/chats/999", nil)
	c.Set("orgID", uint64(100))
	c.Params = []gin.Param{{Key: "chatID", Value: "999"}}

	mockChatService.On("GetByID", uint64(999)).Return((*domain.Chat)(nil), nil)

	handler.GetChat(c)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Chat not found", response["error"])

	mockChatService.AssertExpectations(t)
}

func TestChatHandler_GetChat_InvalidID(t *testing.T) {
	handler, _, _ := setupChatHandler()

	c, w := setupGinContext("GET", "/v1/chats/invalid", nil)
	c.Params = []gin.Param{{Key: "chatID", Value: "invalid"}}

	handler.GetChat(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Invalid chat ID", response["error"])
}

func TestChatHandler_GetChat_WrongOrganization(t *testing.T) {
	handler, mockChatService, _ := setupChatHandler()

	expectedChat := fixtures.CreateTestChat(1)
	expectedChat.ID = 1
	expectedChat.OrganizationID = 200 // Different org

	c, w := setupGinContext("GET", "/v1/chats/1", nil)
	c.Set("orgID", uint64(100))
	c.Params = []gin.Param{{Key: "chatID", Value: "1"}}

	mockChatService.On("GetByID", uint64(1)).Return(expectedChat, nil)

	handler.GetChat(c)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "You do not have permission to access this chat", response["error"])

	mockChatService.AssertExpectations(t)
}

func TestChatHandler_GetChat_WithMessages(t *testing.T) {
	handler, mockChatService, mockMessageService := setupChatHandler()

	expectedChat := fixtures.CreateTestChat(1)
	expectedChat.ID = 1
	expectedChat.OrganizationID = 100

	expectedMessages := []domain.Message{
		*fixtures.CreateTestMessage(1),
	}

	c, w := setupGinContext("GET", "/v1/chats/1?include_messages=true", nil)
	c.Set("orgID", uint64(100))
	c.Params = []gin.Param{{Key: "chatID", Value: "1"}}
	c.Request.URL.RawQuery = "include_messages=true"

	mockChatService.On("GetByID", uint64(1)).Return(expectedChat, nil)
	mockMessageService.On("GetByChatID", uint64(1)).Return(expectedMessages, nil)

	handler.GetChat(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response GetChatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedChat.ID, response.ID)
	assert.Len(t, response.Messages, 1)

	mockChatService.AssertExpectations(t)
	mockMessageService.AssertExpectations(t)
}

func TestChatHandler_ListChats_Success(t *testing.T) {
	handler, mockChatService, _ := setupChatHandler()

	expectedChats := []domain.Chat{
		*fixtures.CreateTestChat(1),
		*fixtures.CreateTestChat(1),
	}

	c, w := setupGinContext("GET", "/v1/chats?limit=10&offset=0", nil)
	c.Set("orgID", uint64(100))
	c.Request.URL.RawQuery = "limit=10&offset=0"

	mockChatService.On("GetByOrganizationID", uint64(100), 10, 0).Return(expectedChats, nil)

	handler.ListChats(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []domain.Chat
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response, 2)

	mockChatService.AssertExpectations(t)
}

func TestChatHandler_UpdateChat_Success(t *testing.T) {
	handler, mockChatService, _ := setupChatHandler()

	existingChat := fixtures.CreateTestChat(1)
	existingChat.ID = 1
	existingChat.OrganizationID = 100

	reqBody := UpdateChatRequest{
		Title: "Updated Title",
		Tags:  []string{"updated", "test"},
	}

	c, w := setupGinContext("PATCH", "/v1/chats/1", reqBody)
	c.Set("orgID", uint64(100))
	c.Params = []gin.Param{{Key: "chatID", Value: "1"}}

	mockChatService.On("GetByID", uint64(1)).Return(existingChat, nil)
	mockChatService.On("UpdateChat", mock.AnythingOfType("*domain.Chat")).Return(nil)

	handler.UpdateChat(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response GetChatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Updated Title", response.Title)
	assert.Equal(t, []string{"updated", "test"}, response.ParsedTags)

	mockChatService.AssertExpectations(t)
}

func TestChatHandler_DeleteChat_Success(t *testing.T) {
	handler, mockChatService, _ := setupChatHandler()

	existingChat := fixtures.CreateTestChat(1)
	existingChat.ID = 1
	existingChat.OrganizationID = 100

	c, w := setupGinContext("DELETE", "/v1/chats/1", nil)
	c.Set("orgID", uint64(100))
	c.Params = []gin.Param{{Key: "chatID", Value: "1"}}

	mockChatService.On("GetByID", uint64(1)).Return(existingChat, nil)
	mockChatService.On("DeleteChat", uint64(1)).Return(nil)

	handler.DeleteChat(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Chat deleted successfully", response["message"])

	mockChatService.AssertExpectations(t)
}

func TestChatHandler_DeleteChat_NotFound(t *testing.T) {
	handler, mockChatService, _ := setupChatHandler()

	c, w := setupGinContext("DELETE", "/v1/chats/999", nil)
	c.Set("orgID", uint64(100))
	c.Params = []gin.Param{{Key: "chatID", Value: "999"}}

	mockChatService.On("GetByID", uint64(999)).Return((*domain.Chat)(nil), nil)

	handler.DeleteChat(c)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Chat not found", response["error"])

	mockChatService.AssertExpectations(t)
}

func TestChatHandler_DeleteChat_WrongOrganization(t *testing.T) {
	handler, mockChatService, _ := setupChatHandler()

	existingChat := fixtures.CreateTestChat(1)
	existingChat.ID = 1
	existingChat.OrganizationID = 200 // Different org

	c, w := setupGinContext("DELETE", "/v1/chats/1", nil)
	c.Set("orgID", uint64(100))
	c.Params = []gin.Param{{Key: "chatID", Value: "1"}}

	mockChatService.On("GetByID", uint64(1)).Return(existingChat, nil)

	handler.DeleteChat(c)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "You do not have permission to delete this chat", response["error"])

	mockChatService.AssertExpectations(t)
}

// Helper function to create uint64 pointer.
func uintPtr(i uint64) *uint64 {
	return &i
}

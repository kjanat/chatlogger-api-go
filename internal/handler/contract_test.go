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

// TestAPIContractValidation tests that API endpoints maintain their expected contracts.
func TestAPIContractValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Chat API Contracts", func(t *testing.T) {
		mockChatService := &mocks.MockChatService{}
		mockMessageService := &mocks.MockMessageService{}

		handler := NewChatHandler(mockChatService, mockMessageService)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("orgID", uint64(1))
			c.Next()
		})

		// Register routes
		router.POST("/chats", handler.CreateChat)
		router.GET("/chats/:chatID", handler.GetChat)
		router.PUT("/chats/:chatID", handler.UpdateChat)
		router.DELETE("/chats/:chatID", handler.DeleteChat)
		router.GET("/chats", handler.ListChats)

		t.Run("POST /chats - Contract Validation", func(t *testing.T) {
			// Setup mock
			mockChatService.On("CreateChat", mock.AnythingOfType("*domain.Chat")).Return(nil)

			// Test valid request
			requestBody := map[string]interface{}{
				"title":    "Test Chat",
				"user_id":  1,
				"tags":     []string{"test", "contract"},
				"metadata": map[string]interface{}{"test": true},
			}

			jsonBody, _ := json.Marshal(requestBody)
			req, _ := http.NewRequest("POST", "/chats", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Contract validation
			assert.Equal(t, http.StatusCreated, w.Code, "Should return 201 Created")
			assert.Equal(
				t,
				"application/json; charset=utf-8",
				w.Header().Get("Content-Type"),
				"Should return JSON content type",
			)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Response should be valid JSON")

			// Validate response structure - CreateChat returns success message and chat_id
			assert.Contains(t, response, "message", "Response should contain success message")
			assert.Contains(t, response, "chat_id", "Response should contain chat_id")

			// Validate data types
			assert.IsType(t, "", response["message"], "Message should be string")
			assert.IsType(t, float64(0), response["chat_id"], "Chat ID should be numeric")
			assert.Equal(
				t,
				"Chat created successfully",
				response["message"],
				"Should return success message",
			)

			mockChatService.AssertExpectations(t)
		})

		t.Run("GET /chats/:id - Contract Validation", func(t *testing.T) {
			// Setup mock
			testChat := fixtures.CreateTestChat(1)
			testChat.ID = 1
			testChat.OrganizationID = 1

			mockChatService.On("GetByID", uint64(1)).Return(testChat, nil)

			req, _ := http.NewRequest("GET", "/chats/1", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Contract validation
			assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK")
			assert.Equal(
				t,
				"application/json; charset=utf-8",
				w.Header().Get("Content-Type"),
				"Should return JSON content type",
			)

			var response GetChatResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Response should be valid JSON")

			// Validate required fields
			assert.NotZero(t, response.ID, "Response should contain valid ID")
			assert.NotEmpty(t, response.Title, "Response should contain title")
			assert.NotZero(t, response.OrganizationID, "Response should contain organization_id")
			assert.NotZero(t, response.CreatedAt, "Response should contain created_at")
			assert.NotZero(t, response.UpdatedAt, "Response should contain updated_at")

			mockChatService.AssertExpectations(t)
		})

		t.Run("Error Response Contract", func(t *testing.T) {
			// Test 404 Not Found contract
			mockChatService.On("GetByID", uint64(999)).Return((*domain.Chat)(nil), nil)

			req, _ := http.NewRequest("GET", "/chats/999", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code, "Should return 404 Not Found")

			var errorResponse map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
			require.NoError(t, err, "Error response should be valid JSON")

			// Validate error response structure
			assert.Contains(t, errorResponse, "error", "Error response should contain error field")
			assert.IsType(t, "", errorResponse["error"], "Error should be string")

			mockChatService.AssertExpectations(t)
		})

		t.Run("PUT /chats/:id - Contract Validation", func(t *testing.T) {
			// Setup mock
			existingChat := fixtures.CreateTestChat(1)
			existingChat.ID = 1
			existingChat.OrganizationID = 1

			mockChatService.On("GetByID", uint64(1)).Return(existingChat, nil)
			mockChatService.On("UpdateChat", mock.AnythingOfType("*domain.Chat")).Return(nil)

			updateBody := map[string]interface{}{
				"title":    "Updated Chat",
				"tags":     []string{"updated"},
				"metadata": map[string]interface{}{"updated": true},
			}

			jsonBody, _ := json.Marshal(updateBody)
			req, _ := http.NewRequest("PUT", "/chats/1", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Contract validation
			assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK")

			var response GetChatResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Response should be valid JSON")

			// Validate response contains updated chat data
			assert.Equal(t, "Updated Chat", response.Title, "Should return updated title")
			assert.NotZero(t, response.ID, "Should contain chat ID")

			mockChatService.AssertExpectations(t)
		})

		t.Run("DELETE /chats/:id - Contract Validation", func(t *testing.T) {
			// Setup mock
			deleteChat := fixtures.CreateTestChat(1)
			deleteChat.ID = 1
			deleteChat.OrganizationID = 1
			mockChatService.On("GetByID", uint64(1)).Return(deleteChat, nil)
			mockChatService.On("DeleteChat", uint64(1)).Return(nil)

			req, _ := http.NewRequest("DELETE", "/chats/1", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Contract validation
			assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK")

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Response should be valid JSON")

			// Validate response contains message
			assert.Contains(t, response, "message", "Response should contain success message")

			mockChatService.AssertExpectations(t)
		})

		t.Run("GET /chats - List Contract Validation", func(t *testing.T) {
			// Setup mock
			chat1 := fixtures.CreateTestChat(1)
			chat1.ID = 1
			chat1.OrganizationID = 1
			chat2 := fixtures.CreateTestChat(1)
			chat2.ID = 2
			chat2.OrganizationID = 1
			chats := []domain.Chat{*chat1, *chat2}

			mockChatService.On("GetByOrganizationID", uint64(1), 20, 0).Return(chats, nil)

			req, _ := http.NewRequest("GET", "/chats", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Contract validation
			assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK")

			var response []domain.Chat
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Response should be valid JSON array")

			// Validate array structure
			assert.Len(t, response, 2, "Should return expected number of chats")

			if len(response) > 0 {
				chat := response[0]
				assert.NotZero(t, chat.ID, "Chat should have ID")
				assert.NotEmpty(t, chat.Title, "Chat should have title")
				assert.NotZero(t, chat.OrganizationID, "Chat should have organization_id")
			}

			mockChatService.AssertExpectations(t)
		})
	})
}

// TestAPIVersionCompatibility tests that API maintains backward compatibility.
func TestAPIVersionCompatibility(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Request Format Compatibility", func(t *testing.T) {
		mockChatService := &mocks.MockChatService{}
		mockMessageService := &mocks.MockMessageService{}

		handler := NewChatHandler(mockChatService, mockMessageService)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("orgID", uint64(1))
			c.Next()
		})
		router.POST("/chats", handler.CreateChat)

		// Test backward compatibility with older request formats
		testCases := []struct {
			name        string
			requestBody map[string]interface{}
			expectValid bool
		}{
			{
				name: "Current format with all fields",
				requestBody: map[string]interface{}{
					"title":    "Test Chat",
					"user_id":  1,
					"tags":     []string{"test"},
					"metadata": map[string]interface{}{"test": true},
				},
				expectValid: true,
			},
			{
				name: "Minimal format (backward compatibility)",
				requestBody: map[string]interface{}{
					"title": "Test Chat",
				},
				expectValid: true,
			},
			{
				name: "Legacy format with string user_id",
				requestBody: map[string]interface{}{
					"title":   "Test Chat",
					"user_id": "1", // String user_id should be rejected
				},
				expectValid: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if tc.expectValid {
					mockChatService.On("CreateChat", mock.AnythingOfType("*domain.Chat")).
						Return(nil).
						Once()
				}

				jsonBody, _ := json.Marshal(tc.requestBody)
				req, _ := http.NewRequest("POST", "/chats", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				if tc.expectValid {
					assert.Equal(
						t,
						http.StatusCreated,
						w.Code,
						"Should accept valid request format",
					)
				} else {
					assert.NotEqual(t, http.StatusCreated, w.Code, "Should reject invalid request format")
				}
			})
		}

		mockChatService.AssertExpectations(t)
	})
}

// TestErrorResponseContracts tests that error responses follow consistent format.
func TestErrorResponseContracts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	errorScenarios := []struct {
		name           string
		setupMock      func(*mocks.MockChatService)
		method         string
		path           string
		body           map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Invalid JSON format",
			setupMock: func(m *mocks.MockChatService) {
				// No mock setup needed for JSON parsing errors
			},
			method:         "POST",
			path:           "/chats",
			body:           nil, // Will send invalid JSON
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request data",
		},
		{
			name: "Chat not found",
			setupMock: func(m *mocks.MockChatService) {
				m.On("GetByID", uint64(999)).Return((*domain.Chat)(nil), nil)
			},
			method:         "GET",
			path:           "/chats/999",
			expectedStatus: http.StatusNotFound,
			expectedError:  "Chat not found",
		},
		{
			name: "Invalid chat ID",
			setupMock: func(m *mocks.MockChatService) {
				// No mock setup needed for parameter validation
			},
			method:         "GET",
			path:           "/chats/invalid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid chat ID",
		},
	}

	for _, scenario := range errorScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			mockChatService := &mocks.MockChatService{}
			mockMessageService := &mocks.MockMessageService{}

			scenario.setupMock(mockChatService)

			handler := NewChatHandler(mockChatService, mockMessageService)
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("orgID", uint64(1))
				c.Next()
			})

			router.POST("/chats", handler.CreateChat)
			router.GET("/chats/:chatID", handler.GetChat)

			var req *http.Request
			switch {
			case scenario.body != nil:
				jsonBody, _ := json.Marshal(scenario.body)
				req, _ = http.NewRequest(scenario.method, scenario.path, bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
			case scenario.method == "POST":
				// Send invalid JSON
				req, _ = http.NewRequest(
					scenario.method,
					scenario.path,
					bytes.NewBufferString("{invalid json"),
				)
				req.Header.Set("Content-Type", "application/json")
			default:
				req, _ = http.NewRequest(scenario.method, scenario.path, nil)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Validate error response contract
			assert.Equal(t, scenario.expectedStatus, w.Code, "Should return expected HTTP status")
			assert.Equal(
				t,
				"application/json; charset=utf-8",
				w.Header().Get("Content-Type"),
				"Should return JSON content type",
			)

			var errorResponse map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
			require.NoError(t, err, "Error response should be valid JSON")

			// Validate error response structure
			assert.Contains(t, errorResponse, "error", "Error response should contain error field")
			assert.IsType(t, "", errorResponse["error"], "Error should be string")

			if scenario.expectedError != "" {
				errorMessage, ok := errorResponse["error"].(string)
				require.True(t, ok, "Error should be string")
				assert.Contains(
					t,
					errorMessage,
					scenario.expectedError,
					"Error message should contain expected text",
				)
			}

			mockChatService.AssertExpectations(t)
		})
	}
}

// TestResponseHeaderContracts tests that response headers are consistent.
func TestResponseHeaderContracts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockChatService := &mocks.MockChatService{}
	mockMessageService := &mocks.MockMessageService{}

	handler := NewChatHandler(mockChatService, mockMessageService)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("orgID", uint64(1))
		c.Next()
	})

	router.GET("/chats/:chatID", handler.GetChat)

	// Setup mock
	testChat := fixtures.CreateTestChat(1)
	testChat.ID = 1
	testChat.OrganizationID = 1
	mockChatService.On("GetByID", uint64(1)).Return(testChat, nil)

	req, _ := http.NewRequest("GET", "/chats/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validate standard headers
	assert.Equal(
		t,
		"application/json; charset=utf-8",
		w.Header().Get("Content-Type"),
		"Should set correct Content-Type",
	)
	// Content-Length may not be set in test mode, just check that response has content
	assert.NotEmpty(t, w.Body.String(), "Response should have content")

	// Validate no sensitive headers are exposed
	assert.Empty(t, w.Header().Get("X-Database-Query"), "Should not expose internal details")
	assert.Empty(t, w.Header().Get("X-Debug-Info"), "Should not expose debug information")

	mockChatService.AssertExpectations(t)
}

// TestPaginationContract tests that pagination follows consistent format.
func TestPaginationContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockChatService := &mocks.MockChatService{}
	mockMessageService := &mocks.MockMessageService{}

	handler := NewChatHandler(mockChatService, mockMessageService)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("orgID", uint64(1))
		c.Next()
	})

	router.GET("/chats", handler.ListChats)

	testCases := []struct {
		name   string
		query  string
		limit  int
		offset int
	}{
		{
			name:   "Default pagination",
			query:  "",
			limit:  20,
			offset: 0,
		},
		{
			name:   "Custom limit",
			query:  "?limit=5",
			limit:  5,
			offset: 0,
		},
		{
			name:   "Custom offset",
			query:  "?offset=20",
			limit:  20,
			offset: 20,
		},
		{
			name:   "Both limit and offset",
			query:  "?limit=15&offset=30",
			limit:  15,
			offset: 30,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock with expected parameters
			chats := []domain.Chat{*fixtures.CreateTestChat(1)}
			mockChatService.On("GetByOrganizationID", uint64(1), tc.limit, tc.offset).
				Return(chats, nil).
				Once()

			req, _ := http.NewRequest("GET", "/chats"+tc.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK")

			var response []domain.Chat
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Response should be valid JSON array")
		})
	}

	mockChatService.AssertExpectations(t)
}

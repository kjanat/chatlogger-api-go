package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/chatlogger-api-go/internal/domain"
)

// MessageHandler handles message-related requests.
type MessageHandler struct {
	messageService domain.MessageService
	chatService    domain.ChatService
}

// NewMessageHandler creates a new message handler.
func NewMessageHandler(
	messageService domain.MessageService,
	chatService domain.ChatService,
) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		chatService:    chatService,
	}
}

// CreateMessageRequest represents the request to create a new message.
type CreateMessageRequest struct {
	Role     domain.MessageRole      `binding:"required" json:"role"`
	Content  string                  `binding:"required" json:"content"`
	Metadata *domain.MessageMetadata `                   json:"metadata,omitempty"` // Use structured type
}

// CreateMessage handles the request to create a new message in a chat.
//
//	@Summary		Create Message
//	@Description	Adds a new message to an existing chat session. Can be called via Public API (API Key).
//	@Tags			Messages
//	@Accept			json
//	@Produce		json
//	@Param			slug	path		string					true	"Organization Slug"
//	@Param			chatID	path		uint64					true	"Chat ID"
//	@Param			request	body		CreateMessageRequest	true	"Message Details (role, content, metadata)"
//	@Success		201		{object}	map[string]interface{}	"message: Message created successfully, message_id: uint64"
//	@Failure		400		{object}	map[string]string		"Invalid chat ID or request data (role, content, metadata validation)"
//	@Failure		401		{object}	map[string]string		"Unauthorized (API Key invalid/missing or Org ID not found)"
//	@Failure		403		{object}	map[string]string		"Forbidden (API Key doesn't match slug or chat doesn't belong to org)"
//	@Failure		404		{object}	map[string]string		"Chat not found"
//	@Failure		500		{object}	map[string]string		"Failed to get chat or create message"
//	@Security		ApiKeyAuth
//	@Router			/v1/orgs/{slug}/chats/{chatID}/messages [post]
func (h *MessageHandler) CreateMessage(c *gin.Context) {
	// Get chat ID from URL
	chatID := c.Param("chatID")
	if chatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chat ID is required"})

		return
	}

	id, err := strconv.ParseUint(chatID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})

		return
	}

	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})

		return
	}

	// Get the chat to validate ownership
	chat, err := h.chatService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat: " + err.Error()})

		return
	}

	if chat == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})

		return
	}

	// Get organization ID from context
	orgIDAny, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})

		return
	}
	orgID := orgIDAny.(uint64)

	// Check if the chat belongs to the organization
	if chat.OrganizationID != orgID {
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "You do not have permission to add messages to this chat"},
		)

		return
	}

	// Create message object
	message := &domain.Message{
		ChatID:    chat.ID,
		Role:      req.Role,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	// Set metadata
	if err := message.SetMetadata(req.Metadata); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to process message metadata: " + err.Error()},
		)

		return
	}

	// Validate the message before creating
	if err := message.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message data: " + err.Error()})

		return
	}

	// Create the message
	if err := h.messageService.CreateMessage(message); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to create message: " + err.Error()},
		)

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Message created successfully",
		"message_id": message.ID,
	})
}

// GetMessageResponse enhances the Message domain model for API responses.
type GetMessageResponse struct {
	*domain.Message
	ParsedMetadata *domain.MessageMetadata `json:"metadata,omitempty"` // Parsed metadata
}

// GetMessages handles the request to get all messages for a chat.
//
//	@Summary		Get Chat Messages
//	@Description	Retrieves all messages associated with a specific chat session.
//	@Tags			Messages
//	@Produce		json
//	@Param			chatID	path		uint64				true	"Chat ID"
//	@Success		200		{array}		GetMessageResponse	"List of messages with parsed metadata"
//	@Failure		400		{object}	map[string]string	"Invalid chat ID"
//	@Failure		401		{object}	map[string]string	"Unauthorized (JWT invalid/missing or Org ID not found)"
//	@Failure		403		{object}	map[string]string	"Permission denied (Chat doesn't belong to user's org)"
//	@Failure		404		{object}	map[string]string	"Chat not found"
//	@Failure		500		{object}	map[string]string	"Failed to get chat or messages"
//	@Security		BearerAuth
//	@Router			/v1/chats/{chatID}/messages [get]
func (h *MessageHandler) GetMessages(c *gin.Context) {
	// Get chat ID from URL
	chatID := c.Param("chatID")
	if chatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chat ID is required"})

		return
	}

	id, err := strconv.ParseUint(chatID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})

		return
	}

	// Get the chat to validate ownership
	chat, err := h.chatService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat: " + err.Error()})

		return
	}

	if chat == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})

		return
	}

	// Get organization ID from context
	orgIDAny, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})

		return
	}
	orgID := orgIDAny.(uint64)

	// Check if the chat belongs to the organization
	if chat.OrganizationID != orgID {
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "You do not have permission to view messages in this chat"},
		)

		return
	}

	// Get messages for the chat
	messages, err := h.messageService.GetByChatID(chat.ID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get messages: " + err.Error()},
		)

		return
	}

	// Prepare response with parsed metadata
	responseMessages := make([]GetMessageResponse, len(messages))
	for i, msg := range messages {
		respMsg := GetMessageResponse{
			Message: &msg,
		}
		metadata, err := msg.GetMetadata()
		if err == nil { // Ignore parsing errors
			respMsg.ParsedMetadata = metadata
		}
		// Nullify raw JSON field in response
		respMsg.Metadata = ""
		responseMessages[i] = respMsg
	}

	c.JSON(http.StatusOK, responseMessages)
}

// GetMessageStats handles the request to get message statistics for an organization.
//
//	@Summary		Get Message Statistics
//	@Description	Retrieves aggregated statistics about messages within a specified date range for the user's organization.
//	@Tags			Analytics
//	@Produce		json
//	@Param			start	query		string					false	"Start date (RFC3339 format, e.g., 2023-01-01T00:00:00Z). Defaults to 30 days ago."
//	@Param			end		query		string					false	"End date (RFC3339 format, e.g., 2023-01-31T23:59:59Z). Defaults to now."
//	@Success		200		{object}	map[string]interface{}	"Aggregated message statistics"
//	@Failure		400		{object}	map[string]string		"Invalid date format"
//	@Failure		401		{object}	map[string]string		"Unauthorized (JWT invalid/missing or Org ID not found)"
//	@Failure		500		{object}	map[string]string		"Failed to get message statistics"
//	@Security		BearerAuth
//	@Router			/v1/analytics/messages [get]
func (h *MessageHandler) GetMessageStats(c *gin.Context) {
	// Get organization ID from context
	orgID, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})

		return
	}

	// Type assert organization ID
	orgIDValue, ok := orgID.(uint64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Organization ID type in context"})
		return
	}

	// Parse date range parameters
	startStr := c.DefaultQuery("start", "")
	endStr := c.DefaultQuery("end", "")

	var start, end time.Time

	var err error

	if startStr == "" {
		// Default to 30 days ago
		start = time.Now().AddDate(0, 0, -30)
	} else {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})

			return
		}
	}

	if endStr == "" {
		// Default to now
		end = time.Now()
	} else {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})

			return
		}
	}

	// Get message statistics
	stats, err := h.messageService.GetMessageStats(orgIDValue, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get message statistics"})

		return
	}

	c.JSON(http.StatusOK, stats)
}

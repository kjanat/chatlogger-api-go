package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/chatlogger-api-go/internal/domain"
)

// ChatHandler handles chat-related requests.
type ChatHandler struct {
	chatService    domain.ChatService
	messageService domain.MessageService
}

// NewChatHandler creates a new chat handler.
func NewChatHandler(
	chatService domain.ChatService,
	messageService domain.MessageService,
) *ChatHandler {
	return &ChatHandler{
		chatService:    chatService,
		messageService: messageService,
	}
}

// CreateChatRequest represents the request to create a new chat.
type CreateChatRequest struct {
	Title    string               `json:"title"`
	Tags     []string             `json:"tags,omitempty"`
	Metadata *domain.ChatMetadata `json:"metadata,omitempty"` // Use the structured type
	UserID   *uint64              `json:"user_id,omitempty"`  // Optional, for anonymous chats
}

// CreateChat handles the request to create a new chat.
//
//	@Summary		Create Chat
//	@Description	Creates a new chat session for the organization. Can be called via Dashboard (JWT) or Public API (API Key).
//	@Tags			Chats
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateChatRequest		true	"Chat Details"
//	@Param			slug	path		string					false	"Organization Slug (Required for Public API)"
//	@Success		201		{object}	map[string]interface{}	"message: Chat created successfully, chat_id: uint64"
//	@Failure		400		{object}	map[string]string		"Invalid request data"
//	@Failure		401		{object}	map[string]string		"Unauthorized (JWT or API Key invalid/missing, or Org ID not found)"
//	@Failure		403		{object}	map[string]string		"Forbidden (API Key doesn't match slug)"
//	@Failure		500		{object}	map[string]string		"Failed to create chat or process tags/metadata"
//	@Security		BearerAuth || ApiKeyAuth
//	@Router			/v1/chats [post]
//	@Router			/v1/orgs/{slug}/chats [post]
func (h *ChatHandler) CreateChat(c *gin.Context) {
	var req CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Get organization ID from context (set by middleware)
	orgIDAny, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return
	}
	orgID := orgIDAny.(uint64)

	// Get a userID if available from the JWT context (may not exist for API key auth)
	var userID *uint64
	userIDInterface, exists := c.Get("userID")
	if exists {
		uid := userIDInterface.(uint64)
		userID = &uid
	} else {
		// Use the userID from the request if provided (primarily for API key scenarios)
		userID = req.UserID
	}

	// Create chat object
	chat := &domain.Chat{
		OrganizationID: orgID,
		UserID:         userID,
		Title:          req.Title,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Set tags using the helper method
	if err := chat.SetTags(req.Tags); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to process tags: " + err.Error()},
		)
		return
	}

	// Set metadata using the helper method
	if err := chat.SetMetadata(req.Metadata); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to process metadata: " + err.Error()},
		)
		return
	}

	// Create the chat
	if err := h.chatService.CreateChat(chat); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to create chat: " + err.Error()},
		)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Chat created successfully",
		"chat_id": chat.ID,
	})
}

// GetChatResponse enhances the Chat domain model for API responses.
type GetChatResponse struct {
	*domain.Chat
	ParsedMetadata *domain.ChatMetadata `json:"metadata,omitempty"` // Parsed metadata
	ParsedTags     []string             `json:"tags,omitempty"`     // Parsed tags
}

// GetChat handles the request to get a chat by ID.
//
//	@Summary		Get Chat by ID
//	@Description	Retrieves details for a specific chat session, optionally including messages.
//	@Tags			Chats
//	@Produce		json
//	@Param			chatID				path		uint64				true	"Chat ID"
//	@Param			include_messages	query		bool				false	"Include messages in the response"	default(false)
//	@Success		200					{object}	GetChatResponse		"Chat details"
//	@Failure		400					{object}	map[string]string	"Invalid chat ID"
//	@Failure		401					{object}	map[string]string	"Unauthorized (JWT invalid/missing or Org ID not found)"
//	@Failure		403					{object}	map[string]string	"Permission denied (Chat doesn't belong to user's org)"
//	@Failure		404					{object}	map[string]string	"Chat not found"
//	@Failure		500					{object}	map[string]string	"Failed to get chat or messages"
//	@Security		BearerAuth
//	@Router			/v1/chats/{chatID} [get]
func (h *ChatHandler) GetChat(c *gin.Context) {
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

	chat, err := h.chatService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat: " + err.Error()})
		return
	}

	if chat == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	orgIDAny, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return
	}
	orgID, ok := orgIDAny.(uint64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Organization ID type in context"})
		return
	}

	if chat.OrganizationID != orgID {
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "You do not have permission to access this chat"},
		)
		return
	}

	// Prepare response object
	response := &GetChatResponse{
		Chat: chat,
	}

	// Parse metadata and tags for the response
	metadata, err := chat.GetMetadata()
	if err == nil { // Ignore parsing errors, return raw if needed
		response.ParsedMetadata = metadata
	}
	tags, err := chat.GetTags()
	if err == nil { // Ignore parsing errors
		response.ParsedTags = tags
	}

	// Optionally include messages
	includeMessages := c.Query("include_messages") == "true"
	if includeMessages {
		messages, err := h.messageService.GetByChatID(chat.ID)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Failed to get messages: " + err.Error()},
			)
			return
		}
		chat.Messages = messages
	}

	// Nullify the raw JSON fields in the response to avoid redundancy
	response.Metadata = ""
	response.Tags = ""

	c.JSON(http.StatusOK, response)
}

// ListChats handles the request to list chats for the current organization.
//
//	@Summary		List Chats
//	@Description	Retrieves a paginated list of chat sessions for the user's organization.
//	@Tags			Chats
//	@Produce		json
//	@Param			limit	query		int					false	"Number of chats per page"	default(20)
//	@Param			offset	query		int					false	"Offset for pagination"		default(0)
//	@Success		200		{array}		domain.Chat			"List of chats"
//	@Failure		401		{object}	map[string]string	"Unauthorized (JWT invalid/missing or Org ID not found)"
//	@Failure		500		{object}	map[string]string	"Failed to list chats"
//	@Security		BearerAuth
//	@Router			/v1/chats [get]
func (h *ChatHandler) ListChats(c *gin.Context) {
	// Get organization ID from context
	orgID, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Get chats
	chats, err := h.chatService.GetByOrganizationID(orgID.(uint64), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list chats"})
		return
	}

	c.JSON(http.StatusOK, chats)
}

// UpdateChatRequest represents the request to update a chat.
type UpdateChatRequest struct {
	Title    string               `json:"title,omitempty"`
	Tags     []string             `json:"tags,omitempty"`
	Metadata *domain.ChatMetadata `json:"metadata,omitempty"` // Use the structured type
}

// UpdateChat handles the request to update a chat.
//
//	@Summary		Update Chat
//	@Description	Updates the title, tags, or metadata of an existing chat session.
//	@Tags			Chats
//	@Accept			json
//	@Produce		json
//	@Param			chatID	path		uint64				true	"Chat ID"
//	@Param			request	body		UpdateChatRequest	true	"Fields to update"
//	@Success		200		{object}	GetChatResponse		"Updated chat details"
//	@Failure		400		{object}	map[string]string	"Invalid chat ID or request data"
//	@Failure		401		{object}	map[string]string	"Unauthorized (JWT invalid/missing or Org ID not found)"
//	@Failure		403		{object}	map[string]string	"Permission denied (Chat doesn't belong to user's org)"
//	@Failure		404		{object}	map[string]string	"Chat not found"
//	@Failure		500		{object}	map[string]string	"Failed to get or update chat, or process tags/metadata"
//	@Security		BearerAuth
//	@Router			/v1/chats/{chatID} [patch]
func (h *ChatHandler) UpdateChat(c *gin.Context) {
	id, req, err := h.parseUpdateChatRequest(c)
	if err != nil {
		return // Error already sent in response
	}

	chat, err := h.validateChatAccess(c, id)
	if err != nil {
		return // Error already sent in response
	}

	if h.applyChatUpdates(c, chat, req) {
		h.saveChatUpdates(c, chat)
	}

	h.sendChatResponse(c, chat)
}

// parseUpdateChatRequest parses and validates the update request.
func (h *ChatHandler) parseUpdateChatRequest(c *gin.Context) (uint64, *UpdateChatRequest, error) {
	chatID := c.Param("chatID")
	if chatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chat ID is required"})
		return 0, nil, fmt.Errorf("missing chat ID")
	}

	id, err := strconv.ParseUint(chatID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return 0, nil, fmt.Errorf("failed to parse chat ID: %w", err)
	}

	var req UpdateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return 0, nil, fmt.Errorf("failed to bind JSON request: %w", err)
	}

	return id, &req, nil
}

// validateChatAccess retrieves the chat and validates user permissions.
func (h *ChatHandler) validateChatAccess(c *gin.Context, id uint64) (*domain.Chat, error) {
	chat, err := h.chatService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat: " + err.Error()})
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}

	if chat == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return nil, fmt.Errorf("chat not found")
	}

	orgIDAny, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return nil, fmt.Errorf("missing org ID")
	}

	orgID := orgIDAny.(uint64)
	if chat.OrganizationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You do not have permission to update this chat",
		})
		return nil, fmt.Errorf("permission denied")
	}

	return chat, nil
}

// applyChatUpdates applies the requested updates to the chat and returns whether any changes were made.
func (h *ChatHandler) applyChatUpdates(
	c *gin.Context,
	chat *domain.Chat,
	req *UpdateChatRequest,
) bool {
	updated := false

	if req.Title != "" {
		chat.Title = req.Title
		updated = true
	}

	if req.Tags != nil {
		if err := chat.SetTags(req.Tags); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to process tags: " + err.Error(),
			})
			return false
		}
		updated = true
	}

	if req.Metadata != nil {
		if err := chat.SetMetadata(req.Metadata); err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Failed to process metadata: " + err.Error()},
			)
			return false
		}
		updated = true
	}

	return updated
}

// saveChatUpdates saves the updated chat to the database.
func (h *ChatHandler) saveChatUpdates(c *gin.Context, chat *domain.Chat) {
	chat.UpdatedAt = time.Now()
	if err := h.chatService.UpdateChat(chat); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to update chat: " + err.Error()},
		)
	}
}

// sendChatResponse prepares and sends the response with updated chat data.
func (h *ChatHandler) sendChatResponse(c *gin.Context, chat *domain.Chat) {
	response := &GetChatResponse{Chat: chat}

	if metadata, _ := chat.GetMetadata(); metadata != nil {
		response.ParsedMetadata = metadata
	}

	if tags, _ := chat.GetTags(); tags != nil {
		response.ParsedTags = tags
	}

	response.Metadata = ""
	response.Tags = ""

	c.JSON(http.StatusOK, response)
}

// DeleteChat handles the request to delete a chat.
//
//	@Summary		Delete Chat
//	@Description	Deletes a chat session and its associated messages.
//	@Tags			Chats
//	@Produce		json
//	@Param			chatID	path		uint64				true	"Chat ID"
//	@Success		200		{object}	map[string]string	"Chat deleted successfully"
//	@Failure		400		{object}	map[string]string	"Invalid chat ID"
//	@Failure		401		{object}	map[string]string	"Unauthorized (JWT invalid/missing or Org ID not found)"
//	@Failure		403		{object}	map[string]string	"Permission denied (Chat doesn't belong to user's org)"
//	@Failure		404		{object}	map[string]string	"Chat not found"
//	@Failure		500		{object}	map[string]string	"Failed to get or delete chat"
//	@Security		BearerAuth
//	@Router			/v1/chats/{chatID} [delete]
func (h *ChatHandler) DeleteChat(c *gin.Context) {
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

	// Get the chat
	chat, err := h.chatService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat"})
		return
	}

	if chat == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	// Get organization ID from context
	orgID, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return
	}

	// Check if the chat belongs to the organization
	if chat.OrganizationID != orgID.(uint64) {
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "You do not have permission to delete this chat"},
		)
		return
	}

	// Delete the chat
	if err := h.chatService.DeleteChat(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted successfully"})
}

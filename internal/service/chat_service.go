package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
)

// ChatService implements the domain.ChatService interface.
type ChatService struct {
	chatRepo domain.ChatRepository
}

// NewChatService creates a new chat service.
func NewChatService(chatRepo domain.ChatRepository) domain.ChatService {
	return &ChatService{
		chatRepo: chatRepo,
	}
}

// CreateChat creates a new chat.
func (s *ChatService) CreateChat(chat *domain.Chat) error {
	// Set timestamps
	chat.CreatedAt = time.Now()
	chat.UpdatedAt = time.Now()

	// Create the chat
	if err := s.chatRepo.Create(chat); err != nil {
		return fmt.Errorf("failed to create chat: %w", err)
	}
	return nil
}

// GetByID gets a chat by ID.
func (s *ChatService) GetByID(id uint64) (*domain.Chat, error) {
	chat, err := s.chatRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find chat by ID: %w", err)
	}
	return chat, nil
}

// GetByOrganizationID gets chats by organization ID with pagination.
func (s *ChatService) GetByOrganizationID(orgID uint64, limit, offset int) ([]domain.Chat, error) {
	chats, err := s.chatRepo.FindByOrganizationID(orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find chats by organization ID: %w", err)
	}
	return chats, nil
}

// GetByUserID gets chats by user ID with pagination.
func (s *ChatService) GetByUserID(userID uint64, limit, offset int) ([]domain.Chat, error) {
	chats, err := s.chatRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find chats by user ID: %w", err)
	}
	return chats, nil
}

// UpdateChat updates a chat.
func (s *ChatService) UpdateChat(chat *domain.Chat) error {
	// Get the existing chat
	existingChat, err := s.chatRepo.FindByID(chat.ID)
	if err != nil {
		return fmt.Errorf("error finding chat: %w", err)
	}

	if existingChat == nil {
		return errors.New("chat not found")
	}

	// Update timestamp
	chat.UpdatedAt = time.Now()

	// Update the chat
	if err := s.chatRepo.Update(chat); err != nil {
		return fmt.Errorf("failed to update chat: %w", err)
	}
	return nil
}

// DeleteChat deletes a chat.
func (s *ChatService) DeleteChat(id uint64) error {
	if err := s.chatRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}
	return nil
}

// GetChatStats gets chat statistics for an organization.
func (s *ChatService) GetChatStats(
	orgID uint64,
	start, end time.Time,
) (map[string]interface{}, error) {
	// Get chat count in date range
	chatCount, err := s.chatRepo.CountByOrgIDAndDateRange(orgID, start, end)
	if err != nil {
		return nil, fmt.Errorf("error getting chat count: %w", err)
	}

	// Get tag statistics
	tagStats, err := s.chatRepo.GetTagStats(orgID)
	if err != nil {
		return nil, fmt.Errorf("error getting tag stats: %w", err)
	}

	// Combine statistics
	stats := map[string]interface{}{
		"total_chats": chatCount,
		"tag_stats":   tagStats,
		"date_range": map[string]string{
			"start": start.Format(time.RFC3339),
			"end":   end.Format(time.RFC3339),
		},
	}

	return stats, nil
}

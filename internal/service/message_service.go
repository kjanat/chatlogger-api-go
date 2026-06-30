package service

import (
	"context"
	"fmt"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
)

// MessageService implements the domain.MessageService interface.
type MessageService struct {
	messageRepo domain.MessageRepository
}

// NewMessageService creates a new message service.
func NewMessageService(messageRepo domain.MessageRepository) domain.MessageService {
	return &MessageService{
		messageRepo: messageRepo,
	}
}

// CreateMessage creates a new message.
func (s *MessageService) CreateMessage(ctx context.Context, message *domain.Message) error {
	// Validate the message
	if err := message.Validate(); err != nil {
		return fmt.Errorf("invalid message: %w", err)
	}

	// Set timestamp
	message.CreatedAt = time.Now()

	// Create the message
	if err := s.messageRepo.Create(ctx, message); err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	return nil
}

// GetByID gets a message by ID.
func (s *MessageService) GetByID(ctx context.Context, id uint64) (*domain.Message, error) {
	message, err := s.messageRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find message by ID: %w", err)
	}
	return message, nil
}

// GetByChatID gets messages by chat ID.
func (s *MessageService) GetByChatID(ctx context.Context, chatID uint64) ([]domain.Message, error) {
	messages, err := s.messageRepo.FindByChatID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to find messages by chat ID: %w", err)
	}
	return messages, nil
}

// GetMessageStats gets message statistics for an organization.
func (s *MessageService) GetMessageStats(
	ctx context.Context,
	orgID uint64,
	start, end time.Time,
) (map[string]interface{}, error) {
	// Get message count in date range
	messageCount, err := s.messageRepo.CountByOrgIDAndDateRange(ctx, orgID, start, end)
	if err != nil {
		return nil, fmt.Errorf("error getting message count: %w", err)
	}

	// Get role statistics
	roleStats, err := s.messageRepo.GetRoleStats(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("error getting role stats: %w", err)
	}

	// Combine statistics
	stats := map[string]interface{}{
		"total_messages": messageCount,
		"by_role":        roleStats,
		"date_range": map[string]string{
			"start": start.Format(time.RFC3339),
			"end":   end.Format(time.RFC3339),
		},
	}

	return stats, nil
}

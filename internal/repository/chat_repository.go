package repository

import (
	"context"
	"errors"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"gorm.io/gorm"
)

// ChatRepo implements the domain.ChatRepository interface.
type ChatRepo struct {
	db *Database
}

// NewChatRepository creates a new chat repository.
func NewChatRepository(db *Database) domain.ChatRepository {
	return &ChatRepo{db: db}
}

// Create creates a new chat.
func (r *ChatRepo) Create(ctx context.Context, chat *domain.Chat) error {
	return r.db.WithContext(ctx).Create(chat).Error
}

// FindByID finds a chat by ID.
func (r *ChatRepo) FindByID(ctx context.Context, id uint64) (*domain.Chat, error) {
	var chat domain.Chat

	err := r.db.WithContext(ctx).First(&chat, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &chat, nil
}

// FindByOrganizationID finds chats by organization ID with pagination.
func (r *ChatRepo) FindByOrganizationID(
	ctx context.Context,
	orgID uint64,
	limit, offset int,
) ([]domain.Chat, error) {
	var chats []domain.Chat
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&chats).
		Error

	return chats, err
}

// FindByUserID finds chats by user ID with pagination.
func (r *ChatRepo) FindByUserID(
	ctx context.Context,
	userID uint64,
	limit, offset int,
) ([]domain.Chat, error) {
	var chats []domain.Chat
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&chats).
		Error

	return chats, err
}

// Update updates a chat.
func (r *ChatRepo) Update(ctx context.Context, chat *domain.Chat) error {
	return r.db.WithContext(ctx).Save(chat).Error
}

// Delete deletes a chat by ID.
func (r *ChatRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.Chat{}, id).Error
}

// CountByOrgIDAndDateRange counts chats in a date range for an organization.
func (r *ChatRepo) CountByOrgIDAndDateRange(
	ctx context.Context,
	orgID uint64,
	start, end time.Time,
) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Chat{}).
		Where("organization_id = ? AND created_at BETWEEN ? AND ?", orgID, start, end).
		Count(&count).Error

	return count, err
}

// GetTagStats gets statistics for tags in an organization.
func (r *ChatRepo) GetTagStats(ctx context.Context, orgID uint64) (map[string]int64, error) {
	// This is a simpler version. For a production implementation,
	// we'd use jsonb functions in PostgreSQL to extract and count tags.
	var chats []domain.Chat

	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&chats).Error
	if err != nil {
		return nil, err
	}

	tagStats := make(map[string]int64)

	// Parse tags from each chat and count occurrences
	for _, chat := range chats {
		tags, err := chat.GetTags()
		if err != nil {
			// Skip chats with invalid tag JSON rather than failing entirely
			continue
		}

		for _, tag := range tags {
			if tag != "" { // Skip empty tags
				tagStats[tag]++
			}
		}
	}

	return tagStats, nil
}

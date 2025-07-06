package repository

import (
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
func (r *ChatRepo) Create(chat *domain.Chat) error {
	return r.db.Create(chat).Error
}

// FindByID finds a chat by ID.
func (r *ChatRepo) FindByID(id uint64) (*domain.Chat, error) {
	var chat domain.Chat

	err := r.db.First(&chat, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &chat, nil
}

// FindByOrganizationID finds chats by organization ID with pagination.
func (r *ChatRepo) FindByOrganizationID(orgID uint64, limit, offset int) ([]domain.Chat, error) {
	var chats []domain.Chat
	err := r.db.Where("organization_id = ?", orgID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&chats).
		Error

	return chats, err
}

// FindByUserID finds chats by user ID with pagination.
func (r *ChatRepo) FindByUserID(userID uint64, limit, offset int) ([]domain.Chat, error) {
	var chats []domain.Chat
	err := r.db.Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&chats).
		Error

	return chats, err
}

// Update updates a chat.
func (r *ChatRepo) Update(chat *domain.Chat) error {
	return r.db.Save(chat).Error
}

// Delete deletes a chat by ID.
func (r *ChatRepo) Delete(id uint64) error {
	return r.db.Delete(&domain.Chat{}, id).Error
}

// CountByOrgIDAndDateRange counts chats in a date range for an organization.
func (r *ChatRepo) CountByOrgIDAndDateRange(orgID uint64, start, end time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Chat{}).
		Where("organization_id = ? AND created_at BETWEEN ? AND ?", orgID, start, end).
		Count(&count).Error

	return count, err
}

// GetTagStats gets statistics for tags in an organization.
func (r *ChatRepo) GetTagStats(orgID uint64) (map[string]int64, error) {
	// This is a simpler version. For a production implementation,
	// we'd use jsonb functions in PostgreSQL to extract and count tags.
	var chats []domain.Chat

	err := r.db.Where("organization_id = ?", orgID).Find(&chats).Error
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

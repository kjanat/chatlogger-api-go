package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

// ChatMetadata represents extended information about a chat session.
type ChatMetadata struct {
	IPAddress        string  `json:"ip_address,omitempty"`
	CountryCode      string  `json:"country_code,omitempty"`  // ISO-3166 country code (e.g., "NL")
	LanguageCode     string  `json:"language_code,omitempty"` // ISO-639 language code (e.g., "tr")
	SessionID        string  `json:"session_id,omitempty"`
	Sentiment        string  `json:"sentiment,omitempty"`
	IsEscalated      bool    `json:"is_escalated,omitempty"`
	IsForwardedToHR  bool    `json:"is_forwarded_to_hr,omitempty"`
	TranscriptLink   string  `json:"transcript_link,omitempty"`
	TokenCount       int     `json:"token_count,omitempty"`       // Total tokens for the chat session
	AvgResponseTime  float64 `json:"avg_response_time,omitempty"` // In seconds
	QuestionCategory string  `json:"question_category,omitempty"`
	UserRating       *int    `json:"user_rating,omitempty"` // Optional user rating
}

// Chat represents a conversation session.
type Chat struct {
	ID             uint64       `gorm:"primaryKey"                json:"id"`
	OrganizationID uint64       `gorm:"not null;index"            json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"-"`
	UserID         *uint64      `                                 json:"user_id,omitempty"` // Nullable for anonymous chats
	User           *User        `gorm:"foreignKey:UserID"         json:"-"`
	Title          string       `gorm:"size:255"                  json:"title"`
	Tags           string       `gorm:"type:jsonb"                json:"tags"`     // JSON array of tags as string
	Metadata       string       `gorm:"type:jsonb"                json:"metadata"` // Store ChatMetadata as JSON string
	CreatedAt      time.Time    `                                 json:"created_at"`
	UpdatedAt      time.Time    `                                 json:"updated_at"`
	Messages       []Message    `gorm:"foreignKey:ChatID"         json:"messages,omitempty"`
}

// GetTags parses the JSON tags string into a slice.
func (c *Chat) GetTags() ([]string, error) {
	var tags []string
	if c.Tags == "" || c.Tags == "null" {
		return []string{}, nil
	}
	err := json.Unmarshal([]byte(c.Tags), &tags)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}
	return tags, nil
}

// SetTags converts a slice of tags into a JSON string.
func (c *Chat) SetTags(tags []string) error {
	if tags == nil {
		tags = []string{} // Ensure empty array instead of null
	}
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}
	c.Tags = string(tagsJSON)
	return nil
}

// GetMetadata parses the JSON metadata string into the ChatMetadata struct.
func (c *Chat) GetMetadata() (*ChatMetadata, error) {
	var metadata ChatMetadata
	if c.Metadata == "" || c.Metadata == "null" {
		return &metadata, nil // Return empty struct if no metadata
	}
	err := json.Unmarshal([]byte(c.Metadata), &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}
	return &metadata, nil
}

// SetMetadata converts the ChatMetadata struct into a JSON string.
func (c *Chat) SetMetadata(metadata *ChatMetadata) error {
	if metadata == nil {
		c.Metadata = "{}" // Store empty JSON object if nil
		return nil
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	c.Metadata = string(metadataJSON)
	return nil
}

// ChatRepository defines the interface for chat data operations.
type ChatRepository interface {
	Create(chat *Chat) error
	FindByID(id uint64) (*Chat, error)
	FindByOrganizationID(orgID uint64, limit, offset int) ([]Chat, error)
	FindByUserID(userID uint64, limit, offset int) ([]Chat, error)
	Update(chat *Chat) error
	Delete(id uint64) error
	CountByOrgIDAndDateRange(orgID uint64, start, end time.Time) (int64, error)
	GetTagStats(orgID uint64) (map[string]int64, error)
}

// ChatService defines the interface for chat business logic.
type ChatService interface {
	CreateChat(chat *Chat) error
	GetByID(id uint64) (*Chat, error)
	GetByOrganizationID(orgID uint64, limit, offset int) ([]Chat, error)
	GetByUserID(userID uint64, limit, offset int) ([]Chat, error)
	UpdateChat(chat *Chat) error
	DeleteChat(id uint64) error
	GetChatStats(orgID uint64, start, end time.Time) (map[string]any, error)
}

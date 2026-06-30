package fixtures

import (
	"fmt"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
)

// OrganizationBuilder provides a fluent interface for building Organization fixtures.
type OrganizationBuilder struct {
	org *domain.Organization
}

// NewOrganizationBuilder creates a new organization builder with defaults.
func NewOrganizationBuilder() *OrganizationBuilder {
	return &OrganizationBuilder{
		org: &domain.Organization{
			Name:      "Test Organization",
			Slug:      "test-org",
			Settings:  `{"test": true}`,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// WithID sets the organization ID.
func (b *OrganizationBuilder) WithID(id uint64) *OrganizationBuilder {
	b.org.ID = id
	return b
}

// WithName sets the organization name.
func (b *OrganizationBuilder) WithName(name string) *OrganizationBuilder {
	b.org.Name = name
	return b
}

// WithSlug sets the organization slug.
func (b *OrganizationBuilder) WithSlug(slug string) *OrganizationBuilder {
	b.org.Slug = slug
	return b
}

// WithSettings sets the organization settings.
func (b *OrganizationBuilder) WithSettings(settings string) *OrganizationBuilder {
	b.org.Settings = settings
	return b
}

// WithTimestamps sets both created and updated timestamps.
func (b *OrganizationBuilder) WithTimestamps(createdAt, updatedAt time.Time) *OrganizationBuilder {
	b.org.CreatedAt = createdAt
	b.org.UpdatedAt = updatedAt
	return b
}

// Build returns the built organization.
func (b *OrganizationBuilder) Build() *domain.Organization {
	return b.org
}

// UserBuilder provides a fluent interface for building User fixtures.
type UserBuilder struct {
	user *domain.User
}

// NewUserBuilder creates a new user builder with defaults.
func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		user: &domain.User{
			OrganizationID: 1,
			Email:          "test@example.com",
			PasswordHash:   "$2a$10$test.hashed.password",
			Role:           domain.RoleUser,
			FirstName:      "Test",
			LastName:       "User",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}
}

// WithID sets the user ID.
func (b *UserBuilder) WithID(id uint64) *UserBuilder {
	b.user.ID = id
	return b
}

// WithOrganizationID sets the organization ID.
func (b *UserBuilder) WithOrganizationID(orgID uint64) *UserBuilder {
	b.user.OrganizationID = orgID
	return b
}

// WithEmail sets the user email.
func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.user.Email = email
	return b
}

// WithRole sets the user role.
func (b *UserBuilder) WithRole(role domain.Role) *UserBuilder {
	b.user.Role = role
	return b
}

// WithName sets both first and last name.
func (b *UserBuilder) WithName(firstName, lastName string) *UserBuilder {
	b.user.FirstName = firstName
	b.user.LastName = lastName
	return b
}

// WithPasswordHash sets the password hash.
func (b *UserBuilder) WithPasswordHash(hash string) *UserBuilder {
	b.user.PasswordHash = hash
	return b
}

// WithTimestamps sets created, updated, and last login timestamps.
func (b *UserBuilder) WithTimestamps(
	createdAt, updatedAt time.Time,
	lastLoginAt *time.Time,
) *UserBuilder {
	b.user.CreatedAt = createdAt
	b.user.UpdatedAt = updatedAt
	b.user.LastLoginAt = lastLoginAt
	return b
}

// AsAdmin configures the user as an admin.
func (b *UserBuilder) AsAdmin() *UserBuilder {
	b.user.Role = domain.RoleAdmin
	b.user.FirstName = "Admin"
	b.user.LastName = "User"
	return b
}

// AsSuperAdmin configures the user as a superadmin.
func (b *UserBuilder) AsSuperAdmin() *UserBuilder {
	b.user.Role = domain.RoleSuperAdmin
	b.user.FirstName = "Super"
	b.user.LastName = "Admin"
	return b
}

// AsViewer configures the user as a viewer.
func (b *UserBuilder) AsViewer() *UserBuilder {
	b.user.Role = domain.RoleViewer
	b.user.FirstName = "Viewer"
	b.user.LastName = "User"
	return b
}

// Build returns the built user.
func (b *UserBuilder) Build() *domain.User {
	return b.user
}

// ChatBuilder provides a fluent interface for building Chat fixtures.
type ChatBuilder struct {
	chat *domain.Chat
}

// NewChatBuilder creates a new chat builder with defaults.
func NewChatBuilder() *ChatBuilder {
	userID := uint64(1)
	return &ChatBuilder{
		chat: &domain.Chat{
			OrganizationID: 1,
			UserID:         &userID,
			Title:          "Test Chat",
			Tags:           `["test", "chat"]`,
			Metadata:       `{"test": "metadata"}`,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}
}

// WithID sets the chat ID.
func (b *ChatBuilder) WithID(id uint64) *ChatBuilder {
	b.chat.ID = id
	return b
}

// WithOrganizationID sets the organization ID.
func (b *ChatBuilder) WithOrganizationID(orgID uint64) *ChatBuilder {
	b.chat.OrganizationID = orgID
	return b
}

// WithUserID sets the user ID.
func (b *ChatBuilder) WithUserID(userID uint64) *ChatBuilder {
	b.chat.UserID = &userID
	return b
}

// WithUser sets the user ID from a user object.
func (b *ChatBuilder) WithUser(user *domain.User) *ChatBuilder {
	if user != nil {
		b.chat.UserID = &user.ID
		b.chat.OrganizationID = user.OrganizationID
	}
	return b
}

// WithTitle sets the chat title.
func (b *ChatBuilder) WithTitle(title string) *ChatBuilder {
	b.chat.Title = title
	return b
}

// WithTags sets the chat tags.
func (b *ChatBuilder) WithTags(tags string) *ChatBuilder {
	b.chat.Tags = tags
	return b
}

// WithMetadata sets the chat metadata.
func (b *ChatBuilder) WithMetadata(metadata string) *ChatBuilder {
	b.chat.Metadata = metadata
	return b
}

// WithTimestamps sets both created and updated timestamps.
func (b *ChatBuilder) WithTimestamps(createdAt, updatedAt time.Time) *ChatBuilder {
	b.chat.CreatedAt = createdAt
	b.chat.UpdatedAt = updatedAt
	return b
}

// Build returns the built chat.
func (b *ChatBuilder) Build() *domain.Chat {
	return b.chat
}

// MessageBuilder provides a fluent interface for building Message fixtures.
type MessageBuilder struct {
	message *domain.Message
}

// NewMessageBuilder creates a new message builder with defaults.
func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{
		message: &domain.Message{
			ChatID:    1,
			Role:      domain.MessageRoleUser,
			Content:   "Test message content",
			Metadata:  `{"token_count": 10}`,
			CreatedAt: time.Now(),
		},
	}
}

// WithID sets the message ID.
func (b *MessageBuilder) WithID(id uint64) *MessageBuilder {
	b.message.ID = id
	return b
}

// WithChatID sets the chat ID.
func (b *MessageBuilder) WithChatID(chatID uint64) *MessageBuilder {
	b.message.ChatID = chatID
	return b
}

// WithChat sets the chat ID from a chat object.
func (b *MessageBuilder) WithChat(chat *domain.Chat) *MessageBuilder {
	if chat != nil {
		b.message.ChatID = chat.ID
	}
	return b
}

// WithRole sets the message role.
func (b *MessageBuilder) WithRole(role domain.MessageRole) *MessageBuilder {
	b.message.Role = role
	return b
}

// WithContent sets the message content.
func (b *MessageBuilder) WithContent(content string) *MessageBuilder {
	b.message.Content = content
	return b
}

// WithMetadata sets the message metadata.
func (b *MessageBuilder) WithMetadata(metadata string) *MessageBuilder {
	b.message.Metadata = metadata
	return b
}

// WithCreatedAt sets the created timestamp.
func (b *MessageBuilder) WithCreatedAt(createdAt time.Time) *MessageBuilder {
	b.message.CreatedAt = createdAt
	return b
}

// AsUserMessage configures the message as a user message.
func (b *MessageBuilder) AsUserMessage() *MessageBuilder {
	b.message.Role = domain.MessageRoleUser
	b.message.Content = "User message content"
	return b
}

// AsAssistantMessage configures the message as an assistant message.
func (b *MessageBuilder) AsAssistantMessage() *MessageBuilder {
	b.message.Role = domain.MessageRoleAssistant
	b.message.Content = "Assistant response content"
	return b
}

// AsSystemMessage configures the message as a system message.
func (b *MessageBuilder) AsSystemMessage() *MessageBuilder {
	b.message.Role = domain.MessageRoleSystem
	b.message.Content = "System message content"
	return b
}

// Build returns the built message.
func (b *MessageBuilder) Build() *domain.Message {
	return b.message
}

// ExportBuilder provides a fluent interface for building Export fixtures.
type ExportBuilder struct {
	export *domain.Export
}

// NewExportBuilder creates a new export builder with defaults.
func NewExportBuilder() *ExportBuilder {
	return &ExportBuilder{
		export: &domain.Export{
			OrganizationID: 1,
			UserID:         1,
			Format:         domain.ExportFormatJSON,
			Type:           domain.ExportTypeAll,
			Status:         domain.ExportStatusPending,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}
}

// WithID sets the export ID.
func (b *ExportBuilder) WithID(id uint64) *ExportBuilder {
	b.export.ID = id
	return b
}

// WithOrganizationID sets the organization ID.
func (b *ExportBuilder) WithOrganizationID(orgID uint64) *ExportBuilder {
	b.export.OrganizationID = orgID
	return b
}

// WithUserID sets the user ID.
func (b *ExportBuilder) WithUserID(userID uint64) *ExportBuilder {
	b.export.UserID = userID
	return b
}

// WithUser sets the user and organization IDs from a user object.
func (b *ExportBuilder) WithUser(user *domain.User) *ExportBuilder {
	if user != nil {
		b.export.UserID = user.ID
		b.export.OrganizationID = user.OrganizationID
	}
	return b
}

// WithFormat sets the export format.
func (b *ExportBuilder) WithFormat(format domain.ExportFormat) *ExportBuilder {
	b.export.Format = format
	return b
}

// WithType sets the export type.
func (b *ExportBuilder) WithType(exportType domain.ExportType) *ExportBuilder {
	b.export.Type = exportType
	return b
}

// WithStatus sets the export status.
func (b *ExportBuilder) WithStatus(status domain.ExportStatus) *ExportBuilder {
	b.export.Status = status
	return b
}

// WithFilePath sets the file path.
func (b *ExportBuilder) WithFilePath(filePath string) *ExportBuilder {
	b.export.FilePath = filePath
	return b
}

// WithError sets the error message.
func (b *ExportBuilder) WithError(error string) *ExportBuilder {
	b.export.Error = error
	return b
}

// WithTimestamps sets both created and updated timestamps.
func (b *ExportBuilder) WithTimestamps(createdAt, updatedAt time.Time) *ExportBuilder {
	b.export.CreatedAt = createdAt
	b.export.UpdatedAt = updatedAt
	return b
}

// AsCSVExport configures the export as CSV format.
func (b *ExportBuilder) AsCSVExport() *ExportBuilder {
	b.export.Format = domain.ExportFormatCSV
	return b
}

// AsJSONExport configures the export as JSON format.
func (b *ExportBuilder) AsJSONExport() *ExportBuilder {
	b.export.Format = domain.ExportFormatJSON
	return b
}

// AsCompleted configures the export as completed.
func (b *ExportBuilder) AsCompleted() *ExportBuilder {
	b.export.Status = domain.ExportStatusCompleted
	b.export.FilePath = "/tmp/export.json"
	return b
}

// AsFailed configures the export as failed.
func (b *ExportBuilder) AsFailed() *ExportBuilder {
	b.export.Status = domain.ExportStatusFailed
	b.export.Error = "Export failed"
	return b
}

// Build returns the built export.
func (b *ExportBuilder) Build() *domain.Export {
	return b.export
}

// APIKeyBuilder provides a fluent interface for building APIKey fixtures.
type APIKeyBuilder struct {
	apiKey *domain.APIKey
}

// NewAPIKeyBuilder creates a new API key builder with defaults.
func NewAPIKeyBuilder() *APIKeyBuilder {
	return &APIKeyBuilder{
		apiKey: &domain.APIKey{
			OrganizationID: 1,
			HashedKey:      "$2a$10$hashed.api.key",
			Label:          "Test API Key",
			CreatedAt:      time.Now(),
		},
	}
}

// WithID sets the API key ID.
func (b *APIKeyBuilder) WithID(id uint64) *APIKeyBuilder {
	b.apiKey.ID = id
	return b
}

// WithOrganizationID sets the organization ID.
func (b *APIKeyBuilder) WithOrganizationID(orgID uint64) *APIKeyBuilder {
	b.apiKey.OrganizationID = orgID
	return b
}

// WithHashedKey sets the hashed key.
func (b *APIKeyBuilder) WithHashedKey(hashedKey string) *APIKeyBuilder {
	b.apiKey.HashedKey = hashedKey
	return b
}

// WithLabel sets the API key label.
func (b *APIKeyBuilder) WithLabel(label string) *APIKeyBuilder {
	b.apiKey.Label = label
	return b
}

// WithCreatedAt sets the created timestamp.
func (b *APIKeyBuilder) WithCreatedAt(createdAt time.Time) *APIKeyBuilder {
	b.apiKey.CreatedAt = createdAt
	return b
}

// Build returns the built API key.
func (b *APIKeyBuilder) Build() *domain.APIKey {
	return b.apiKey
}

// TestDataFactory provides factory methods for creating related test data.
type TestDataFactory struct{}

// NewTestDataFactory creates a new test data factory.
func NewTestDataFactory() *TestDataFactory {
	return &TestDataFactory{}
}

// CreateOrganizationWithUsers creates an organization with multiple users.
func (f *TestDataFactory) CreateOrganizationWithUsers(
	orgName string,
	userCount int,
) (*domain.Organization, []*domain.User) {
	org := NewOrganizationBuilder().
		WithName(orgName).
		WithSlug(fmt.Sprintf("%s-slug", orgName)).
		Build()

	var users []*domain.User
	for i := 0; i < userCount; i++ {
		user := NewUserBuilder().
			WithOrganizationID(org.ID).
			WithEmail(fmt.Sprintf("user%d@%s.com", i+1, orgName)).
			WithName(fmt.Sprintf("User%d", i+1), "Test").
			Build()
		users = append(users, user)
	}

	return org, users
}

// CreateChatWithMessages creates a chat with multiple messages.
func (f *TestDataFactory) CreateChatWithMessages(
	user *domain.User,
	title string,
	messageCount int,
) (*domain.Chat, []*domain.Message) {
	chat := NewChatBuilder().
		WithUser(user).
		WithTitle(title).
		Build()

	var messages []*domain.Message
	for i := 0; i < messageCount; i++ {
		role := domain.MessageRoleUser
		if i%2 == 1 {
			role = domain.MessageRoleAssistant
		}

		message := NewMessageBuilder().
			WithChat(chat).
			WithRole(role).
			WithContent(fmt.Sprintf("Message %d content", i+1)).
			Build()
		messages = append(messages, message)
	}

	return chat, messages
}

// CreateConversation creates a full conversation scenario.
func (f *TestDataFactory) CreateConversation(
	orgName string,
) (*domain.Organization, *domain.User, *domain.Chat, []*domain.Message) {
	org := NewOrganizationBuilder().
		WithName(orgName).
		WithSlug(fmt.Sprintf("%s-slug", orgName)).
		Build()

	user := NewUserBuilder().
		WithOrganizationID(org.ID).
		WithEmail(fmt.Sprintf("user@%s.com", orgName)).
		Build()

	chat, messages := f.CreateChatWithMessages(user, "Test Conversation", 4)

	return org, user, chat, messages
}

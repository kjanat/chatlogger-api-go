package fixtures

import (
	"fmt"
	"testing"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/stretchr/testify/assert"
)

// Example tests demonstrating the new fixture builder pattern

func TestOrganizationBuilder_Example(t *testing.T) {
	// Using the new builder pattern for flexible organization creation
	org := NewOrganizationBuilder().
		WithName("Tech Startup").
		WithSlug("tech-startup").
		WithSettings(`{"theme": "dark", "notifications": true}`).
		Build()

	assert.Equal(t, "Tech Startup", org.Name)
	assert.Equal(t, "tech-startup", org.Slug)
	assert.Contains(t, org.Settings, "dark")
	assert.NotZero(t, org.CreatedAt)
	assert.NotZero(t, org.UpdatedAt)
}

func TestUserBuilder_Example(t *testing.T) {
	// Creating different types of users with the builder pattern
	org := NewOrganizationBuilder().WithID(1).Build()

	// Regular user
	user := NewUserBuilder().
		WithOrganizationID(org.ID).
		WithEmail("john.doe@example.com").
		WithName("John", "Doe").
		Build()

	assert.Equal(t, org.ID, user.OrganizationID)
	assert.Equal(t, "john.doe@example.com", user.Email)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, domain.RoleUser, user.Role)

	// Admin user using convenience method
	admin := NewUserBuilder().
		WithOrganizationID(org.ID).
		WithEmail("admin@example.com").
		AsAdmin().
		Build()

	assert.Equal(t, domain.RoleAdmin, admin.Role)
	assert.Equal(t, "Admin", admin.FirstName)

	// Superadmin user
	superadmin := NewUserBuilder().
		WithOrganizationID(org.ID).
		WithEmail("superadmin@example.com").
		AsSuperAdmin().
		Build()

	assert.Equal(t, domain.RoleSuperAdmin, superadmin.Role)
	assert.Equal(t, "Super", superadmin.FirstName)
}

func TestChatBuilder_Example(t *testing.T) {
	// Creating related data using builders
	user := NewUserBuilder().
		WithID(1).
		WithOrganizationID(1).
		Build()

	// Chat with automatic relationship handling
	chat := NewChatBuilder().
		WithUser(user). // Automatically sets user ID and org ID
		WithTitle("Project Discussion").
		WithTags(`["project", "discussion", "team"]`).
		WithMetadata(`{"priority": "high", "category": "work"}`).
		Build()

	assert.Equal(t, user.ID, *chat.UserID)
	assert.Equal(t, user.OrganizationID, chat.OrganizationID)
	assert.Equal(t, "Project Discussion", chat.Title)
	assert.Contains(t, chat.Tags, "project")
	assert.Contains(t, chat.Metadata, "high")
}

func TestMessageBuilder_Example(t *testing.T) {
	// Creating a chat with messages
	user := NewUserBuilder().WithID(1).Build()
	chat := NewChatBuilder().WithUser(user).WithID(1).Build()

	// Different types of messages
	userMessage := NewMessageBuilder().
		WithChat(chat).
		AsUserMessage().
		WithContent("Hello, I need help with this feature").
		Build()

	assistantMessage := NewMessageBuilder().
		WithChat(chat).
		AsAssistantMessage().
		WithContent("I'd be happy to help! What specific part are you struggling with?").
		Build()

	systemMessage := NewMessageBuilder().
		WithChat(chat).
		AsSystemMessage().
		WithContent("User john.doe@example.com joined the conversation").
		Build()

	assert.Equal(t, chat.ID, userMessage.ChatID)
	assert.Equal(t, domain.MessageRoleUser, userMessage.Role)
	assert.Contains(t, userMessage.Content, "help")

	assert.Equal(t, domain.MessageRoleAssistant, assistantMessage.Role)
	assert.Contains(t, assistantMessage.Content, "happy")

	assert.Equal(t, domain.MessageRoleSystem, systemMessage.Role)
	assert.Contains(t, systemMessage.Content, "joined")
}

func TestExportBuilder_Example(t *testing.T) {
	user := NewUserBuilder().WithID(1).WithOrganizationID(1).Build()

	// Different export scenarios
	pendingExport := NewExportBuilder().
		WithUser(user).
		AsJSONExport().
		Build()

	completedExport := NewExportBuilder().
		WithUser(user).
		AsCSVExport().
		AsCompleted().
		Build()

	failedExport := NewExportBuilder().
		WithUser(user).
		AsFailed().
		Build()

	assert.Equal(t, domain.ExportStatusPending, pendingExport.Status)
	assert.Equal(t, domain.ExportFormatJSON, pendingExport.Format)

	assert.Equal(t, domain.ExportStatusCompleted, completedExport.Status)
	assert.Equal(t, domain.ExportFormatCSV, completedExport.Format)
	assert.NotEmpty(t, completedExport.FilePath)

	assert.Equal(t, domain.ExportStatusFailed, failedExport.Status)
	assert.Equal(t, "Export failed", failedExport.Error)
}

func TestTestDataFactory_Example(t *testing.T) {
	factory := NewTestDataFactory()

	// Create organization with multiple users
	org, users := factory.CreateOrganizationWithUsers("TestCorp", 3)

	assert.Equal(t, "TestCorp", org.Name)
	assert.Equal(t, "TestCorp-slug", org.Slug)
	assert.Len(t, users, 3)

	for i, user := range users {
		assert.Equal(t, org.ID, user.OrganizationID)
		assert.Contains(t, user.Email, "TestCorp")
		assert.Equal(t, fmt.Sprintf("User%d", i+1), user.FirstName)
	}

	// Create chat with messages
	user := users[0]
	chat, messages := factory.CreateChatWithMessages(user, "Feature Discussion", 6)

	assert.Equal(t, user.ID, *chat.UserID)
	assert.Equal(t, user.OrganizationID, chat.OrganizationID)
	assert.Equal(t, "Feature Discussion", chat.Title)
	assert.Len(t, messages, 6)

	// Verify alternating message roles
	for i, message := range messages {
		assert.Equal(t, chat.ID, message.ChatID)
		if i%2 == 0 {
			assert.Equal(t, domain.MessageRoleUser, message.Role)
		} else {
			assert.Equal(t, domain.MessageRoleAssistant, message.Role)
		}
	}
}

func TestTestDataFactory_CompleteScenario(t *testing.T) {
	factory := NewTestDataFactory()

	// Create a complete conversation scenario
	org, user, chat, messages := factory.CreateConversation("SupportChat")

	// Verify the complete scenario
	assert.Equal(t, "SupportChat", org.Name)
	assert.Equal(t, "SupportChat-slug", org.Slug)

	assert.Equal(t, org.ID, user.OrganizationID)
	assert.Contains(t, user.Email, "SupportChat")

	assert.Equal(t, user.ID, *chat.UserID)
	assert.Equal(t, org.ID, chat.OrganizationID)
	assert.Equal(t, "Test Conversation", chat.Title)

	assert.Len(t, messages, 4)
	for _, message := range messages {
		assert.Equal(t, chat.ID, message.ChatID)
		assert.NotEmpty(t, message.Content)
	}
}

func TestBuilderPattern_Comparison(t *testing.T) {
	// Demonstrating the difference between old and new patterns

	t.Run("Old Pattern", func(t *testing.T) {
		// Old pattern - less flexible, hard-coded values
		chat := CreateTestChat(1)
		chat.OrganizationID = 1 // Manual override needed

		assert.Equal(t, "Test Chat", chat.Title)       // Fixed title
		assert.Equal(t, `["test", "chat"]`, chat.Tags) // Fixed tags
	})

	t.Run("New Builder Pattern", func(t *testing.T) {
		// New pattern - flexible, explicit, reusable
		user := NewUserBuilder().WithID(1).WithOrganizationID(1).Build()

		chat := NewChatBuilder().
			WithUser(user). // Automatic relationship handling
			WithTitle("Custom Support Chat").
			WithTags(`["support", "urgent", "customer"]`).
			WithMetadata(`{"customer_tier": "premium"}`).
			Build()

		assert.Equal(t, "Custom Support Chat", chat.Title)
		assert.Contains(t, chat.Tags, "urgent")
		assert.Contains(t, chat.Metadata, "premium")
		assert.Equal(t, user.ID, *chat.UserID)
		assert.Equal(t, user.OrganizationID, chat.OrganizationID)
	})
}

func TestBuilderPattern_Chaining(t *testing.T) {
	// Demonstrating advanced chaining and relationship handling
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	org := NewOrganizationBuilder().
		WithID(1).
		WithName("AdvancedCorp").
		WithTimestamps(yesterday, now).
		Build()

	admin := NewUserBuilder().
		WithID(1).
		WithOrganizationID(org.ID).
		WithEmail("admin@advancedcorp.com").
		AsAdmin().
		WithTimestamps(yesterday, now, &now).
		Build()

	chat := NewChatBuilder().
		WithID(1).
		WithUser(admin).
		WithTitle("Admin Discussion").
		WithTimestamps(yesterday, now).
		Build()

	message := NewMessageBuilder().
		WithID(1).
		WithChat(chat).
		AsUserMessage().
		WithContent("This is an admin message").
		WithCreatedAt(now).
		Build()

	// Verify all relationships are correctly established
	assert.Equal(t, admin.OrganizationID, org.ID)
	assert.Equal(t, *chat.UserID, admin.ID)
	assert.Equal(t, chat.OrganizationID, org.ID)
	assert.Equal(t, message.ChatID, chat.ID)

	// Verify timestamps are properly set
	assert.Equal(t, yesterday.Unix(), org.CreatedAt.Unix())
	assert.Equal(t, now.Unix(), org.UpdatedAt.Unix())
	assert.NotNil(t, admin.LastLoginAt)
	assert.Equal(t, now.Unix(), admin.LastLoginAt.Unix())
}

// Benchmark to show the performance characteristics.
func BenchmarkFixtureCreation(b *testing.B) {
	b.Run("OldPattern", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			org := CreateTestOrganization()
			user := CreateTestUser(org.ID)
			chat := CreateTestChat(user.ID)
			chat.OrganizationID = org.ID
			_ = CreateTestMessage(chat.ID)
		}
	})

	b.Run("NewBuilderPattern", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			org := NewOrganizationBuilder().WithID(1).Build()
			user := NewUserBuilder().WithID(1).WithOrganizationID(org.ID).Build()
			chat := NewChatBuilder().WithID(1).WithUser(user).Build()
			_ = NewMessageBuilder().WithID(1).WithChat(chat).Build()
		}
	})

	b.Run("FactoryPattern", func(b *testing.B) {
		factory := NewTestDataFactory()
		for i := 0; i < b.N; i++ {
			_, _, _, _ = factory.CreateConversation("BenchTest")
		}
	})
}

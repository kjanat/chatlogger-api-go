package fixtures

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
)

// CreateTestOrganization creates a test organization
func CreateTestOrganization() *domain.Organization {
	random := rand.Intn(100000)
	return &domain.Organization{
		Name:      fmt.Sprintf("Test Organization %d", random),
		Slug:      fmt.Sprintf("test-org-%d", random),
		Settings:  `{"test": true}`,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestUser creates a test user
func CreateTestUser(orgID uint64) *domain.User {
	random := rand.Intn(100000)
	return &domain.User{
		OrganizationID: orgID,
		Email:          fmt.Sprintf("test%d@example.com", random),
		PasswordHash:   "$2a$10$test.hashed.password",
		Role:           domain.RoleUser,
		FirstName:      "Test",
		LastName:       "User",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// CreateTestChat creates a test chat
func CreateTestChat(userID uint64) *domain.Chat {
	userIDPtr := &userID
	return &domain.Chat{
		OrganizationID: 1,
		UserID:         userIDPtr,
		Title:          "Test Chat",
		Tags:           `["test", "chat"]`,
		Metadata:       `{"test": "metadata"}`,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// CreateTestMessage creates a test message
func CreateTestMessage(chatID uint64) *domain.Message {
	return &domain.Message{
		ChatID:    chatID,
		Role:      domain.MessageRoleUser,
		Content:   "Test message content",
		Metadata:  `{"token_count": 10}`,
		CreatedAt: time.Now(),
	}
}

// CreateTestExport creates a test export
func CreateTestExport(userID uint64) *domain.Export {
	return &domain.Export{
		OrganizationID: 1,
		UserID:         userID,
		Format:         domain.ExportFormatJSON,
		Type:           domain.ExportTypeAll,
		Status:         domain.ExportStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// CreateTestAPIKey creates a test API key
func CreateTestAPIKey(orgID uint64) *domain.APIKey {
	return &domain.APIKey{
		OrganizationID: orgID,
		HashedKey:      "$2a$10$hashed.api.key",
		Label:          "Test API Key",
		CreatedAt:      time.Now(),
	}
}

// CreateTestSuperadmin creates a test superadmin user
func CreateTestSuperadmin(orgID uint64) *domain.User {
	random := rand.Intn(100000)
	return &domain.User{
		OrganizationID: orgID,
		Email:          fmt.Sprintf("superadmin%d@example.com", random),
		PasswordHash:   "$2a$10$test.hashed.password",
		Role:           domain.RoleSuperAdmin,
		FirstName:      "Super",
		LastName:       "Admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
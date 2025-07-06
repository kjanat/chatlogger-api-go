package mocks

import (
	"fmt"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockChatService is a mock implementation of ChatService.
type MockChatService struct {
	mock.Mock
}

func (m *MockChatService) CreateChat(chat *domain.Chat) error {
	args := m.Called(chat)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock create chat error: %w", err)
	}
	return nil
}

func (m *MockChatService) GetByID(id uint64) (*domain.Chat, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("mock get chat by ID error: %w", err)
		}
		return nil, nil
	}
	chat := args.Get(0).(*domain.Chat)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock get chat by ID error: %w", err)
	}
	return chat, nil
}

func (m *MockChatService) GetByOrganizationID(
	orgID uint64,
	limit, offset int,
) ([]domain.Chat, error) {
	args := m.Called(orgID, limit, offset)
	chats := args.Get(0).([]domain.Chat)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock get chats by organization ID error: %w", err)
	}
	return chats, nil
}

func (m *MockChatService) GetByUserID(userID uint64, limit, offset int) ([]domain.Chat, error) {
	args := m.Called(userID, limit, offset)
	chats := args.Get(0).([]domain.Chat)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock get chats by user ID error: %w", err)
	}
	return chats, nil
}

func (m *MockChatService) UpdateChat(chat *domain.Chat) error {
	args := m.Called(chat)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock update chat error: %w", err)
	}
	return nil
}

func (m *MockChatService) DeleteChat(id uint64) error {
	args := m.Called(id)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock delete chat error: %w", err)
	}
	return nil
}

func (m *MockChatService) GetChatStats(
	orgID uint64,
	start, end time.Time,
) (map[string]interface{}, error) {
	args := m.Called(orgID, start, end)
	stats := args.Get(0).(map[string]interface{})
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock get chat stats error: %w", err)
	}
	return stats, nil
}

// MockMessageService is a mock implementation of MessageService.
type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) CreateMessage(message *domain.Message) error {
	args := m.Called(message)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock create message error: %w", err)
	}
	return nil
}

func (m *MockMessageService) GetByID(id uint64) (*domain.Message, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("mock get message by ID error: %w", err)
		}
		return nil, nil
	}
	message := args.Get(0).(*domain.Message)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock get message by ID error: %w", err)
	}
	return message, nil
}

func (m *MockMessageService) GetByChatID(chatID uint64) ([]domain.Message, error) {
	args := m.Called(chatID)
	messages := args.Get(0).([]domain.Message)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock get messages by chat ID error: %w", err)
	}
	return messages, nil
}

func (m *MockMessageService) GetMessageStats(
	orgID uint64,
	start, end time.Time,
) (map[string]interface{}, error) {
	args := m.Called(orgID, start, end)
	stats := args.Get(0).(map[string]interface{})
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock get message stats error: %w", err)
	}
	return stats, nil
}

// MockOrganizationService is a mock implementation of OrganizationService.
type MockOrganizationService struct {
	mock.Mock
}

func (m *MockOrganizationService) Create(org *domain.Organization) error {
	args := m.Called(org)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock create organization error: %w", err)
	}
	return nil
}

func (m *MockOrganizationService) GetByID(id uint64) (*domain.Organization, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("mock get organization by ID error: %w", err)
		}
		return nil, nil
	}
	org := args.Get(0).(*domain.Organization)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock get organization by ID error: %w", err)
	}
	return org, nil
}

func (m *MockOrganizationService) GetBySlug(slug string) (*domain.Organization, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("mock get organization by slug error: %w", err)
		}
		return nil, nil
	}
	org := args.Get(0).(*domain.Organization)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock get organization by slug error: %w", err)
	}
	return org, nil
}

func (m *MockOrganizationService) Update(org *domain.Organization) error {
	args := m.Called(org)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock update organization error: %w", err)
	}
	return nil
}

func (m *MockOrganizationService) Delete(id uint64) error {
	args := m.Called(id)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock delete organization error: %w", err)
	}
	return nil
}

func (m *MockOrganizationService) List(limit, offset int) ([]domain.Organization, error) {
	args := m.Called(limit, offset)
	orgs := args.Get(0).([]domain.Organization)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock list organizations error: %w", err)
	}
	return orgs, nil
}

// MockAPIKeyService is a mock implementation of APIKeyService.
type MockAPIKeyService struct {
	mock.Mock
}

func (m *MockAPIKeyService) GenerateKey(orgID uint64, label string) (string, error) {
	args := m.Called(orgID, label)
	key := args.String(0)
	if err := args.Error(1); err != nil {
		return "", fmt.Errorf("mock generate key error: %w", err)
	}
	return key, nil
}

func (m *MockAPIKeyService) ValidateKey(rawKey string) (*domain.APIKey, error) {
	args := m.Called(rawKey)
	if args.Get(0) == nil {
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("mock validate key error: %w", err)
		}
		return nil, nil
	}
	apiKey := args.Get(0).(*domain.APIKey)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock validate key error: %w", err)
	}
	return apiKey, nil
}

func (m *MockAPIKeyService) GetByID(id uint64) (*domain.APIKey, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("mock get API key by ID error: %w", err)
		}
		return nil, nil
	}
	apiKey := args.Get(0).(*domain.APIKey)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock get API key by ID error: %w", err)
	}
	return apiKey, nil
}

func (m *MockAPIKeyService) ListByOrganizationID(orgID uint64) ([]domain.APIKey, error) {
	args := m.Called(orgID)
	keys := args.Get(0).([]domain.APIKey)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock list API keys by organization ID error: %w", err)
	}
	return keys, nil
}

func (m *MockAPIKeyService) RevokeKey(id uint64) error {
	args := m.Called(id)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock revoke key error: %w", err)
	}
	return nil
}

func (m *MockAPIKeyService) DeleteKey(id uint64) error {
	args := m.Called(id)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock delete key error: %w", err)
	}
	return nil
}

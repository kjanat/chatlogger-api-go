package mocks

import (
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockChatService is a mock implementation of ChatService
type MockChatService struct {
	mock.Mock
}

func (m *MockChatService) CreateChat(chat *domain.Chat) error {
	args := m.Called(chat)
	return args.Error(0)
}

func (m *MockChatService) GetByID(id uint64) (*domain.Chat, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Chat), args.Error(1)
}

func (m *MockChatService) GetByOrganizationID(orgID uint64, limit, offset int) ([]domain.Chat, error) {
	args := m.Called(orgID, limit, offset)
	return args.Get(0).([]domain.Chat), args.Error(1)
}

func (m *MockChatService) GetByUserID(userID uint64, limit, offset int) ([]domain.Chat, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]domain.Chat), args.Error(1)
}

func (m *MockChatService) UpdateChat(chat *domain.Chat) error {
	args := m.Called(chat)
	return args.Error(0)
}

func (m *MockChatService) DeleteChat(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockChatService) GetChatStats(orgID uint64, start, end time.Time) (map[string]interface{}, error) {
	args := m.Called(orgID, start, end)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// MockMessageService is a mock implementation of MessageService
type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) CreateMessage(message *domain.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockMessageService) GetByID(id uint64) (*domain.Message, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageService) GetByChatID(chatID uint64) ([]domain.Message, error) {
	args := m.Called(chatID)
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (m *MockMessageService) GetMessageStats(orgID uint64, start, end time.Time) (map[string]interface{}, error) {
	args := m.Called(orgID, start, end)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// MockOrganizationService is a mock implementation of OrganizationService
type MockOrganizationService struct {
	mock.Mock
}

func (m *MockOrganizationService) Create(org *domain.Organization) error {
	args := m.Called(org)
	return args.Error(0)
}

func (m *MockOrganizationService) GetByID(id uint64) (*domain.Organization, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrganizationService) GetBySlug(slug string) (*domain.Organization, error) {
	args := m.Called(slug)
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrganizationService) Update(org *domain.Organization) error {
	args := m.Called(org)
	return args.Error(0)
}

func (m *MockOrganizationService) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockOrganizationService) List(limit, offset int) ([]domain.Organization, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]domain.Organization), args.Error(1)
}

// MockAPIKeyService is a mock implementation of APIKeyService
type MockAPIKeyService struct {
	mock.Mock
}

func (m *MockAPIKeyService) GenerateKey(orgID uint64, label string) (string, error) {
	args := m.Called(orgID, label)
	return args.String(0), args.Error(1)
}

func (m *MockAPIKeyService) ValidateKey(rawKey string) (*domain.APIKey, error) {
	args := m.Called(rawKey)
	return args.Get(0).(*domain.APIKey), args.Error(1)
}

func (m *MockAPIKeyService) GetByID(id uint64) (*domain.APIKey, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.APIKey), args.Error(1)
}

func (m *MockAPIKeyService) ListByOrganizationID(orgID uint64) ([]domain.APIKey, error) {
	args := m.Called(orgID)
	return args.Get(0).([]domain.APIKey), args.Error(1)
}

func (m *MockAPIKeyService) RevokeKey(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAPIKeyService) DeleteKey(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}
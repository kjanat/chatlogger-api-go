package mocks

import (
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockChatRepository is a mock implementation of ChatRepository
type MockChatRepository struct {
	mock.Mock
}

func (m *MockChatRepository) Create(chat *domain.Chat) error {
	args := m.Called(chat)
	return args.Error(0)
}

func (m *MockChatRepository) FindByID(id uint64) (*domain.Chat, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Chat), args.Error(1)
}

func (m *MockChatRepository) FindByOrganizationID(orgID uint64, limit, offset int) ([]domain.Chat, error) {
	args := m.Called(orgID, limit, offset)
	return args.Get(0).([]domain.Chat), args.Error(1)
}

func (m *MockChatRepository) FindByUserID(userID uint64, limit, offset int) ([]domain.Chat, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]domain.Chat), args.Error(1)
}

func (m *MockChatRepository) Update(chat *domain.Chat) error {
	args := m.Called(chat)
	return args.Error(0)
}

func (m *MockChatRepository) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockChatRepository) CountByOrgIDAndDateRange(orgID uint64, start, end time.Time) (int64, error) {
	args := m.Called(orgID, start, end)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChatRepository) GetTagStats(orgID uint64) (map[string]int64, error) {
	args := m.Called(orgID)
	return args.Get(0).(map[string]int64), args.Error(1)
}

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uint64) (*domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByOrganizationID(orgID uint64, limit, offset int) ([]domain.User, error) {
	args := m.Called(orgID, limit, offset)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockOrganizationRepository is a mock implementation of OrganizationRepository
type MockOrganizationRepository struct {
	mock.Mock
}

func (m *MockOrganizationRepository) Create(org *domain.Organization) error {
	args := m.Called(org)
	return args.Error(0)
}

func (m *MockOrganizationRepository) FindByID(id uint64) (*domain.Organization, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) FindBySlug(slug string) (*domain.Organization, error) {
	args := m.Called(slug)
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) Update(org *domain.Organization) error {
	args := m.Called(org)
	return args.Error(0)
}

func (m *MockOrganizationRepository) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockOrganizationRepository) List(limit, offset int) ([]domain.Organization, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]domain.Organization), args.Error(1)
}

// MockMessageRepository is a mock implementation of MessageRepository
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Create(message *domain.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockMessageRepository) FindByID(id uint64) (*domain.Message, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageRepository) FindByChatID(chatID uint64) ([]domain.Message, error) {
	args := m.Called(chatID)
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (m *MockMessageRepository) CountByOrgIDAndDateRange(orgID uint64, start, end time.Time) (int64, error) {
	args := m.Called(orgID, start, end)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMessageRepository) GetRoleStats(orgID uint64) (map[domain.MessageRole]int64, error) {
	args := m.Called(orgID)
	return args.Get(0).(map[domain.MessageRole]int64), args.Error(1)
}

// MockExportRepository is a mock implementation of ExportRepository
type MockExportRepository struct {
	mock.Mock
}

func (m *MockExportRepository) Create(export *domain.Export) error {
	args := m.Called(export)
	return args.Error(0)
}

func (m *MockExportRepository) GetByID(id uint64) (*domain.Export, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Export), args.Error(1)
}

func (m *MockExportRepository) GetByOrganizationID(organizationID uint64, limit, offset int) ([]*domain.Export, error) {
	args := m.Called(organizationID, limit, offset)
	return args.Get(0).([]*domain.Export), args.Error(1)
}

func (m *MockExportRepository) UpdateStatus(id uint64, status domain.ExportStatus, errorMsg string) error {
	args := m.Called(id, status, errorMsg)
	return args.Error(0)
}

func (m *MockExportRepository) UpdateFilePath(id uint64, filePath string) error {
	args := m.Called(id, filePath)
	return args.Error(0)
}
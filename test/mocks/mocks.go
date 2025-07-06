package mocks

import (
	"fmt"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockChatRepository is a mock implementation of ChatRepository.
type MockChatRepository struct {
	mock.Mock
}

func (m *MockChatRepository) Create(chat *domain.Chat) error {
	args := m.Called(chat)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock create chat error: %w", err)
	}
	return nil
}

func (m *MockChatRepository) FindByID(id uint64) (*domain.Chat, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("mock find chat by ID error: %w", err)
		}
		return nil, nil
	}
	chat := args.Get(0).(*domain.Chat)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock find chat by ID error: %w", err)
	}
	return chat, nil
}

func (m *MockChatRepository) FindByOrganizationID(
	orgID uint64,
	limit, offset int,
) ([]domain.Chat, error) {
	args := m.Called(orgID, limit, offset)
	chats := args.Get(0).([]domain.Chat)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock find chats by organization ID error: %w", err)
	}
	return chats, nil
}

func (m *MockChatRepository) FindByUserID(userID uint64, limit, offset int) ([]domain.Chat, error) {
	args := m.Called(userID, limit, offset)
	chats := args.Get(0).([]domain.Chat)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock find chats by user ID error: %w", err)
	}
	return chats, nil
}

func (m *MockChatRepository) Update(chat *domain.Chat) error {
	args := m.Called(chat)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock update chat error: %w", err)
	}
	return nil
}

func (m *MockChatRepository) Delete(id uint64) error {
	args := m.Called(id)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock delete chat error: %w", err)
	}
	return nil
}

func (m *MockChatRepository) CountByOrgIDAndDateRange(
	orgID uint64,
	start, end time.Time,
) (int64, error) {
	args := m.Called(orgID, start, end)
	count := args.Get(0).(int64)
	if err := args.Error(1); err != nil {
		return 0, fmt.Errorf("mock count by org ID and date range error: %w", err)
	}
	return count, nil
}

func (m *MockChatRepository) GetTagStats(orgID uint64) (map[string]int64, error) {
	args := m.Called(orgID)
	stats := args.Get(0).(map[string]int64)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock get tag stats error: %w", err)
	}
	return stats, nil
}

// MockUserRepository is a mock implementation of UserRepository.
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock create user error: %w", err)
	}
	return nil
}

func (m *MockUserRepository) FindByID(id uint64) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("mock find user by ID error: %w", err)
		}
		return nil, nil
	}
	user := args.Get(0).(*domain.User)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock find user by ID error: %w", err)
	}
	return user, nil
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("mock find user by email error: %w", err)
		}
		return nil, nil
	}
	user := args.Get(0).(*domain.User)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock find user by email error: %w", err)
	}
	return user, nil
}

func (m *MockUserRepository) FindByOrganizationID(
	orgID uint64,
	limit, offset int,
) ([]domain.User, error) {
	args := m.Called(orgID, limit, offset)
	users := args.Get(0).([]domain.User)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock find users by organization ID error: %w", err)
	}
	return users, nil
}

func (m *MockUserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock update user error: %w", err)
	}
	return nil
}

func (m *MockUserRepository) Delete(id uint64) error {
	args := m.Called(id)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock delete user error: %w", err)
	}
	return nil
}

// MockOrganizationRepository is a mock implementation of OrganizationRepository.
type MockOrganizationRepository struct {
	mock.Mock
}

func (m *MockOrganizationRepository) Create(org *domain.Organization) error {
	args := m.Called(org)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock create organization error: %w", err)
	}
	return nil
}

func (m *MockOrganizationRepository) FindByID(id uint64) (*domain.Organization, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("mock find organization by ID error: %w", err)
		}
		return nil, nil
	}
	org := args.Get(0).(*domain.Organization)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock find organization by ID error: %w", err)
	}
	return org, nil
}

func (m *MockOrganizationRepository) FindBySlug(slug string) (*domain.Organization, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("mock find organization by slug error: %w", err)
		}
		return nil, nil
	}
	org := args.Get(0).(*domain.Organization)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock find organization by slug error: %w", err)
	}
	return org, nil
}

func (m *MockOrganizationRepository) Update(org *domain.Organization) error {
	args := m.Called(org)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock update organization error: %w", err)
	}
	return nil
}

func (m *MockOrganizationRepository) Delete(id uint64) error {
	args := m.Called(id)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock delete organization error: %w", err)
	}
	return nil
}

func (m *MockOrganizationRepository) List(limit, offset int) ([]domain.Organization, error) {
	args := m.Called(limit, offset)
	orgs := args.Get(0).([]domain.Organization)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock list organizations error: %w", err)
	}
	return orgs, nil
}

// MockMessageRepository is a mock implementation of MessageRepository.
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Create(message *domain.Message) error {
	args := m.Called(message)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock create message error: %w", err)
	}
	return nil
}

func (m *MockMessageRepository) FindByID(id uint64) (*domain.Message, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("mock find message by ID error: %w", err)
		}
		return nil, nil
	}
	message := args.Get(0).(*domain.Message)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock find message by ID error: %w", err)
	}
	return message, nil
}

func (m *MockMessageRepository) FindByChatID(chatID uint64) ([]domain.Message, error) {
	args := m.Called(chatID)
	messages := args.Get(0).([]domain.Message)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock find messages by chat ID error: %w", err)
	}
	return messages, nil
}

func (m *MockMessageRepository) CountByOrgIDAndDateRange(
	orgID uint64,
	start, end time.Time,
) (int64, error) {
	args := m.Called(orgID, start, end)
	count := args.Get(0).(int64)
	if err := args.Error(1); err != nil {
		return 0, fmt.Errorf("mock count messages by org ID and date range error: %w", err)
	}
	return count, nil
}

func (m *MockMessageRepository) GetRoleStats(orgID uint64) (map[domain.MessageRole]int64, error) {
	args := m.Called(orgID)
	stats := args.Get(0).(map[domain.MessageRole]int64)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock get role stats error: %w", err)
	}
	return stats, nil
}

// MockExportRepository is a mock implementation of ExportRepository.
type MockExportRepository struct {
	mock.Mock
}

func (m *MockExportRepository) Create(export *domain.Export) error {
	args := m.Called(export)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock create export error: %w", err)
	}
	return nil
}

func (m *MockExportRepository) GetByID(id uint64) (*domain.Export, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("mock get export by ID error: %w", err)
		}
		return nil, nil
	}
	export := args.Get(0).(*domain.Export)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock get export by ID error: %w", err)
	}
	return export, nil
}

func (m *MockExportRepository) GetByOrganizationID(
	organizationID uint64,
	limit, offset int,
) ([]*domain.Export, error) {
	args := m.Called(organizationID, limit, offset)
	exports := args.Get(0).([]*domain.Export)
	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock get exports by organization ID error: %w", err)
	}
	return exports, nil
}

func (m *MockExportRepository) UpdateStatus(
	id uint64,
	status domain.ExportStatus,
	errorMsg string,
) error {
	args := m.Called(id, status, errorMsg)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock update export status error: %w", err)
	}
	return nil
}

func (m *MockExportRepository) UpdateFilePath(id uint64, filePath string) error {
	args := m.Called(id, filePath)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock update export file path error: %w", err)
	}
	return nil
}

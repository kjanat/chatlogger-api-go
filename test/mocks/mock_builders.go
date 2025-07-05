package mocks

import (
	"errors"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/kjanat/chatlogger-api-go/test/fixtures"
	"github.com/stretchr/testify/mock"
)

// MockChatRepositoryBuilder provides a fluent interface for building MockChatRepository
type MockChatRepositoryBuilder struct {
	mock *MockChatRepository
}

// NewMockChatRepositoryBuilder creates a new builder instance
func NewMockChatRepositoryBuilder() *MockChatRepositoryBuilder {
	return &MockChatRepositoryBuilder{
		mock: &MockChatRepository{},
	}
}

// WithCreateReturns configures the Create method mock
func (b *MockChatRepositoryBuilder) WithCreateReturns(err error) *MockChatRepositoryBuilder {
	b.mock.On("Create", mock.AnythingOfType("*domain.Chat")).Return(err)
	return b
}

// WithFindByIDReturns configures the FindByID method mock
func (b *MockChatRepositoryBuilder) WithFindByIDReturns(chat *domain.Chat, err error) *MockChatRepositoryBuilder {
	b.mock.On("FindByID", mock.AnythingOfType("uint64")).Return(chat, err)
	return b
}

// WithFindByOrganizationIDReturns configures the FindByOrganizationID method mock
func (b *MockChatRepositoryBuilder) WithFindByOrganizationIDReturns(chats []domain.Chat, err error) *MockChatRepositoryBuilder {
	b.mock.On("FindByOrganizationID", mock.AnythingOfType("uint64"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(chats, err)
	return b
}

// WithUpdateReturns configures the Update method mock
func (b *MockChatRepositoryBuilder) WithUpdateReturns(err error) *MockChatRepositoryBuilder {
	b.mock.On("Update", mock.AnythingOfType("*domain.Chat")).Return(err)
	return b
}

// WithDeleteReturns configures the Delete method mock
func (b *MockChatRepositoryBuilder) WithDeleteReturns(err error) *MockChatRepositoryBuilder {
	b.mock.On("Delete", mock.AnythingOfType("uint64")).Return(err)
	return b
}

// WithCountByOrgIDAndDateRangeReturns configures the CountByOrgIDAndDateRange method mock
func (b *MockChatRepositoryBuilder) WithCountByOrgIDAndDateRangeReturns(count int64, err error) *MockChatRepositoryBuilder {
	b.mock.On("CountByOrgIDAndDateRange", mock.AnythingOfType("uint64"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(count, err)
	return b
}

// WithGetTagStatsReturns configures the GetTagStats method mock
func (b *MockChatRepositoryBuilder) WithGetTagStatsReturns(stats map[string]int64, err error) *MockChatRepositoryBuilder {
	b.mock.On("GetTagStats", mock.AnythingOfType("uint64")).Return(stats, err)
	return b
}

// Build returns the configured mock
func (b *MockChatRepositoryBuilder) Build() *MockChatRepository {
	return b.mock
}

// MockUserRepositoryBuilder provides a fluent interface for building MockUserRepository
type MockUserRepositoryBuilder struct {
	mock *MockUserRepository
}

// NewMockUserRepositoryBuilder creates a new builder instance
func NewMockUserRepositoryBuilder() *MockUserRepositoryBuilder {
	return &MockUserRepositoryBuilder{
		mock: &MockUserRepository{},
	}
}

// WithCreateReturns configures the Create method mock
func (b *MockUserRepositoryBuilder) WithCreateReturns(err error) *MockUserRepositoryBuilder {
	b.mock.On("Create", mock.AnythingOfType("*domain.User")).Return(err)
	return b
}

// WithFindByIDReturns configures the FindByID method mock
func (b *MockUserRepositoryBuilder) WithFindByIDReturns(user *domain.User, err error) *MockUserRepositoryBuilder {
	b.mock.On("FindByID", mock.AnythingOfType("uint64")).Return(user, err)
	return b
}

// WithFindByEmailReturns configures the FindByEmail method mock
func (b *MockUserRepositoryBuilder) WithFindByEmailReturns(user *domain.User, err error) *MockUserRepositoryBuilder {
	b.mock.On("FindByEmail", mock.AnythingOfType("string")).Return(user, err)
	return b
}

// WithFindByOrganizationIDReturns configures the FindByOrganizationID method mock
func (b *MockUserRepositoryBuilder) WithFindByOrganizationIDReturns(users []domain.User, err error) *MockUserRepositoryBuilder {
	b.mock.On("FindByOrganizationID", mock.AnythingOfType("uint64"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(users, err)
	return b
}

// WithUpdateReturns configures the Update method mock
func (b *MockUserRepositoryBuilder) WithUpdateReturns(err error) *MockUserRepositoryBuilder {
	b.mock.On("Update", mock.AnythingOfType("*domain.User")).Return(err)
	return b
}

// WithDeleteReturns configures the Delete method mock
func (b *MockUserRepositoryBuilder) WithDeleteReturns(err error) *MockUserRepositoryBuilder {
	b.mock.On("Delete", mock.AnythingOfType("uint64")).Return(err)
	return b
}

// Build returns the configured mock
func (b *MockUserRepositoryBuilder) Build() *MockUserRepository {
	return b.mock
}

// MockChatServiceBuilder provides a fluent interface for building MockChatService
type MockChatServiceBuilder struct {
	mock *MockChatService
}

// NewMockChatServiceBuilder creates a new builder instance
func NewMockChatServiceBuilder() *MockChatServiceBuilder {
	return &MockChatServiceBuilder{
		mock: &MockChatService{},
	}
}

// WithCreateChatReturns configures the CreateChat method mock
func (b *MockChatServiceBuilder) WithCreateChatReturns(err error) *MockChatServiceBuilder {
	b.mock.On("CreateChat", mock.AnythingOfType("*domain.Chat")).Return(err)
	return b
}

// WithGetByIDReturns configures the GetByID method mock
func (b *MockChatServiceBuilder) WithGetByIDReturns(chat *domain.Chat, err error) *MockChatServiceBuilder {
	b.mock.On("GetByID", mock.AnythingOfType("uint64")).Return(chat, err)
	return b
}

// WithGetByOrganizationIDReturns configures the GetByOrganizationID method mock
func (b *MockChatServiceBuilder) WithGetByOrganizationIDReturns(chats []domain.Chat, err error) *MockChatServiceBuilder {
	b.mock.On("GetByOrganizationID", mock.AnythingOfType("uint64"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(chats, err)
	return b
}

// WithUpdateChatReturns configures the UpdateChat method mock
func (b *MockChatServiceBuilder) WithUpdateChatReturns(err error) *MockChatServiceBuilder {
	b.mock.On("UpdateChat", mock.AnythingOfType("*domain.Chat")).Return(err)
	return b
}

// WithDeleteChatReturns configures the DeleteChat method mock
func (b *MockChatServiceBuilder) WithDeleteChatReturns(err error) *MockChatServiceBuilder {
	b.mock.On("DeleteChat", mock.AnythingOfType("uint64")).Return(err)
	return b
}

// WithGetChatStatsReturns configures the GetChatStats method mock
func (b *MockChatServiceBuilder) WithGetChatStatsReturns(stats map[string]interface{}, err error) *MockChatServiceBuilder {
	b.mock.On("GetChatStats", mock.AnythingOfType("uint64"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(stats, err)
	return b
}

// Build returns the configured mock
func (b *MockChatServiceBuilder) Build() *MockChatService {
	return b.mock
}

// ServiceTestSetup provides a standardized setup for service tests with common test data
type ServiceTestSetup struct {
	Org     *domain.Organization
	User    *domain.User
	Chat    *domain.Chat
	Message *domain.Message
}

// NewServiceTestSetup creates standard test data for service tests
func NewServiceTestSetup() *ServiceTestSetup {
	org := fixtures.CreateTestOrganization()
	org.ID = 1

	user := fixtures.CreateTestUser(org.ID)
	user.ID = 1

	chat := fixtures.CreateTestChat(user.ID)
	chat.ID = 1
	chat.OrganizationID = org.ID

	message := fixtures.CreateTestMessage(chat.ID)
	message.ID = 1

	return &ServiceTestSetup{
		Org:     org,
		User:    user,
		Chat:    chat,
		Message: message,
	}
}

// DefaultChatServiceMock creates a chat service mock with common successful responses
func DefaultChatServiceMock(setup *ServiceTestSetup) *MockChatService {
	return NewMockChatServiceBuilder().
		WithCreateChatReturns(nil).
		WithGetByIDReturns(setup.Chat, nil).
		WithGetByOrganizationIDReturns([]domain.Chat{*setup.Chat}, nil).
		WithUpdateChatReturns(nil).
		WithDeleteChatReturns(nil).
		WithGetChatStatsReturns(map[string]interface{}{"count": 1}, nil).
		Build()
}

// DefaultChatRepositoryMock creates a chat repository mock with common successful responses
func DefaultChatRepositoryMock(setup *ServiceTestSetup) *MockChatRepository {
	return NewMockChatRepositoryBuilder().
		WithCreateReturns(nil).
		WithFindByIDReturns(setup.Chat, nil).
		WithFindByOrganizationIDReturns([]domain.Chat{*setup.Chat}, nil).
		WithUpdateReturns(nil).
		WithDeleteReturns(nil).
		WithCountByOrgIDAndDateRangeReturns(1, nil).
		WithGetTagStatsReturns(map[string]int64{"test": 1}, nil).
		Build()
}

// DefaultUserRepositoryMock creates a user repository mock with common successful responses
func DefaultUserRepositoryMock(setup *ServiceTestSetup) *MockUserRepository {
	return NewMockUserRepositoryBuilder().
		WithCreateReturns(nil).
		WithFindByIDReturns(setup.User, nil).
		WithFindByEmailReturns(setup.User, nil).
		WithFindByOrganizationIDReturns([]domain.User{*setup.User}, nil).
		WithUpdateReturns(nil).
		WithDeleteReturns(nil).
		Build()
}

// ScenarioBuilder provides methods for creating common test scenarios
type ScenarioBuilder struct {
	setup *ServiceTestSetup
}

// NewScenarioBuilder creates a new scenario builder
func NewScenarioBuilder() *ScenarioBuilder {
	return &ScenarioBuilder{
		setup: NewServiceTestSetup(),
	}
}

// GetSetup returns the test setup data
func (s *ScenarioBuilder) GetSetup() *ServiceTestSetup {
	return s.setup
}

// ChatNotFoundScenario creates mocks for a chat not found scenario
func (s *ScenarioBuilder) ChatNotFoundScenario() *MockChatRepository {
	return NewMockChatRepositoryBuilder().
		WithFindByIDReturns(nil, nil).
		Build()
}

// UserNotFoundScenario creates mocks for a user not found scenario
func (s *ScenarioBuilder) UserNotFoundScenario() *MockUserRepository {
	return NewMockUserRepositoryBuilder().
		WithFindByEmailReturns(nil, nil).
		Build()
}

// DatabaseErrorScenario creates mocks that return database errors
func (s *ScenarioBuilder) DatabaseErrorScenario() (*MockChatRepository, *MockUserRepository) {
	dbErr := errors.New("database error")

	chatRepo := NewMockChatRepositoryBuilder().
		WithCreateReturns(dbErr).
		WithFindByIDReturns(nil, dbErr).
		Build()

	userRepo := NewMockUserRepositoryBuilder().
		WithCreateReturns(dbErr).
		WithFindByIDReturns(nil, dbErr).
		Build()

	return chatRepo, userRepo
}

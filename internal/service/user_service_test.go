package service

import (
	"errors"
	"testing"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/kjanat/chatlogger-api-go/internal/hash"
	"github.com/kjanat/chatlogger-api-go/test/fixtures"
	"github.com/kjanat/chatlogger-api-go/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_Authenticate_Success(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := NewUserService(mockRepo, "test-secret")

	password := "password123"
	hashedPassword, err := hash.GeneratePasswordHash(password, 10)
	require.NoError(t, err)

	user := fixtures.CreateTestUser(1)
	user.Email = "test@example.com"
	user.PasswordHash = hashedPassword

	mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)

	resultUser, token, err := service.Authenticate("test@example.com", password)

	assert.NoError(t, err)
	assert.NotNil(t, resultUser)
	assert.NotEmpty(t, token)
	assert.Equal(t, user.Email, resultUser.Email)
	assert.NotNil(t, resultUser.LastLoginAt)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Authenticate_UserNotFound(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := NewUserService(mockRepo, "test-secret")

	mockRepo.On("FindByEmail", "nonexistent@example.com").Return((*domain.User)(nil), nil)

	resultUser, token, err := service.Authenticate("nonexistent@example.com", "password")

	assert.Error(t, err)
	assert.Nil(t, resultUser)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "invalid email or password")
	mockRepo.AssertExpectations(t)
}

func TestUserService_Authenticate_WrongPassword(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := NewUserService(mockRepo, "test-secret")

	hashedPassword, err := hash.GeneratePasswordHash("correctpassword", 10)
	require.NoError(t, err)

	user := fixtures.CreateTestUser(1)
	user.Email = "test@example.com"
	user.PasswordHash = hashedPassword

	mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)

	resultUser, token, err := service.Authenticate("test@example.com", "wrongpassword")

	assert.Error(t, err)
	assert.Nil(t, resultUser)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "invalid email or password")
	mockRepo.AssertExpectations(t)
}

func TestUserService_Authenticate_FindError(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := NewUserService(mockRepo, "test-secret")

	expectedError := errors.New("database error")
	mockRepo.On("FindByEmail", "test@example.com").Return((*domain.User)(nil), expectedError)

	resultUser, token, err := service.Authenticate("test@example.com", "password")

	assert.Error(t, err)
	assert.Nil(t, resultUser)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "error finding user")
	mockRepo.AssertExpectations(t)
}

func TestUserService_Register_Success(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := NewUserService(mockRepo, "test-secret")

	user := &domain.User{
		OrganizationID: 1,
		Email:          "new@example.com",
		Role:           domain.RoleUser,
		FirstName:      "New",
		LastName:       "User",
	}

	mockRepo.On("FindByEmail", "new@example.com").Return((*domain.User)(nil), nil)
	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	err := service.Register(user, "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, user.PasswordHash)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Register_EmailExists(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := NewUserService(mockRepo, "test-secret")

	existingUser := fixtures.CreateTestUser(1)
	existingUser.Email = "existing@example.com"

	user := &domain.User{
		OrganizationID: 1,
		Email:          "existing@example.com",
		Role:           domain.RoleUser,
	}

	mockRepo.On("FindByEmail", "existing@example.com").Return(existingUser, nil)

	err := service.Register(user, "password123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user with this email already exists")
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByID(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := NewUserService(mockRepo, "test-secret")

	expectedUser := fixtures.CreateTestUser(1)
	expectedUser.ID = 1

	mockRepo.On("FindByID", uint64(1)).Return(expectedUser, nil)

	user, err := service.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByEmail(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := NewUserService(mockRepo, "test-secret")

	expectedUser := fixtures.CreateTestUser(1)
	expectedUser.Email = "test@example.com"

	mockRepo.On("FindByEmail", "test@example.com").Return(expectedUser, nil)

	user, err := service.GetByEmail("test@example.com")

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByOrganizationID(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := NewUserService(mockRepo, "test-secret")

	expectedUsers := []domain.User{
		*fixtures.CreateTestUser(1),
		*fixtures.CreateTestUser(1),
	}

	mockRepo.On("FindByOrganizationID", uint64(1), 10, 0).Return(expectedUsers, nil)

	users, err := service.GetByOrganizationID(1, 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := NewUserService(mockRepo, "test-secret")

	user := fixtures.CreateTestUser(1)
	user.FirstName = "Updated"

	mockRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)

	err := service.UpdateUser(user)

	assert.NoError(t, err)
	assert.NotZero(t, user.UpdatedAt)
	mockRepo.AssertExpectations(t)
}

func TestUserService_ChangePassword_Success(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := NewUserService(mockRepo, "test-secret")

	currentPassword := "currentpass"
	hashedPassword, err := hash.GeneratePasswordHash(currentPassword, 10)
	require.NoError(t, err)

	user := fixtures.CreateTestUser(1)
	user.ID = 1
	user.PasswordHash = hashedPassword

	mockRepo.On("FindByID", uint64(1)).Return(user, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)

	err = service.ChangePassword(1, currentPassword, "newpassword")

	assert.NoError(t, err)
	assert.NotEqual(t, hashedPassword, user.PasswordHash) // Password should be changed
	mockRepo.AssertExpectations(t)
}

func TestUserService_ChangePassword_UserNotFound(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := NewUserService(mockRepo, "test-secret")

	mockRepo.On("FindByID", uint64(999)).Return((*domain.User)(nil), nil)

	err := service.ChangePassword(999, "currentpass", "newpass")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	mockRepo.AssertExpectations(t)
}

func TestUserService_ChangePassword_WrongCurrentPassword(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := NewUserService(mockRepo, "test-secret")

	hashedPassword, err := hash.GeneratePasswordHash("correctpass", 10)
	require.NoError(t, err)

	user := fixtures.CreateTestUser(1)
	user.ID = 1
	user.PasswordHash = hashedPassword

	mockRepo.On("FindByID", uint64(1)).Return(user, nil)

	err = service.ChangePassword(1, "wrongpass", "newpass")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "current password is incorrect")
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := NewUserService(mockRepo, "test-secret")

	mockRepo.On("Delete", uint64(1)).Return(nil)

	err := service.DeleteUser(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser_Error(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := NewUserService(mockRepo, "test-secret")

	expectedError := errors.New("delete error")
	mockRepo.On("Delete", uint64(1)).Return(expectedError)

	err := service.DeleteUser(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete user")
	assert.Contains(t, err.Error(), "delete error")
	mockRepo.AssertExpectations(t)
}

func TestGenerateJWT(t *testing.T) {
	user := &domain.User{
		ID:             1,
		Email:          "test@example.com",
		OrganizationID: 100,
		Role:           domain.RoleUser,
	}

	secret := "test-secret"

	token, err := generateJWT(user, secret)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Token should be a valid JWT (contains two dots)
	assert.Equal(t, 2, countOccurrences(token, "."))
}

func TestJWTClaims_Structure(t *testing.T) {
	claims := &JWTClaims{
		UserID:         1,
		Email:          "test@example.com",
		OrganizationID: 100,
		Role:           domain.RoleAdmin,
	}

	assert.Equal(t, uint64(1), claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, uint64(100), claims.OrganizationID)
	assert.Equal(t, domain.RoleAdmin, claims.Role)
}

// Helper function to count occurrences of a substring.
func countOccurrences(s, substr string) int {
	count := 0
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			count++
		}
	}
	return count
}

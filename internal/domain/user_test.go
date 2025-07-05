package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRole_Constants(t *testing.T) {
	assert.Equal(t, Role("superadmin"), RoleSuperAdmin)
	assert.Equal(t, Role("admin"), RoleAdmin)
	assert.Equal(t, Role("user"), RoleUser)
	assert.Equal(t, Role("viewer"), RoleViewer)
}

func TestUser_Struct(t *testing.T) {
	now := time.Now()
	lastLogin := now.Add(-24 * time.Hour)
	
	user := &User{
		ID:             1,
		OrganizationID: 100,
		Email:          "test@example.com",
		PasswordHash:   "$2a$10$hashedpassword",
		Role:           RoleUser,
		FirstName:      "John",
		LastName:       "Doe",
		CreatedAt:      now,
		UpdatedAt:      now,
		LastLoginAt:    &lastLogin,
	}
	
	assert.Equal(t, uint64(1), user.ID)
	assert.Equal(t, uint64(100), user.OrganizationID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "$2a$10$hashedpassword", user.PasswordHash)
	assert.Equal(t, RoleUser, user.Role)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, now, user.CreatedAt)
	assert.Equal(t, now, user.UpdatedAt)
	assert.NotNil(t, user.LastLoginAt)
	assert.Equal(t, lastLogin, *user.LastLoginAt)
}

func TestUser_NilLastLoginAt(t *testing.T) {
	user := &User{
		ID:             1,
		OrganizationID: 100,
		Email:          "test@example.com",
		PasswordHash:   "$2a$10$hashedpassword",
		Role:           RoleAdmin,
		FirstName:      "Jane",
		LastName:       "Smith",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		LastLoginAt:    nil,
	}
	
	assert.Nil(t, user.LastLoginAt)
}

func TestUser_AllRoles(t *testing.T) {
	roles := []Role{RoleSuperAdmin, RoleAdmin, RoleUser, RoleViewer}
	
	for _, role := range roles {
		user := &User{
			ID:             1,
			OrganizationID: 100,
			Email:          "test@example.com",
			PasswordHash:   "$2a$10$hashedpassword",
			Role:           role,
			FirstName:      "Test",
			LastName:       "User",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		
		assert.Equal(t, role, user.Role)
	}
}

func TestUser_EmptyFields(t *testing.T) {
	user := &User{
		ID:             1,
		OrganizationID: 100,
		Email:          "minimal@example.com",
		PasswordHash:   "$2a$10$hashedpassword",
		Role:           RoleUser,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	
	assert.Empty(t, user.FirstName)
	assert.Empty(t, user.LastName)
	assert.Nil(t, user.LastLoginAt)
	assert.Empty(t, user.Chats)
}
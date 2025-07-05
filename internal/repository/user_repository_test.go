package repository

import (
	"testing"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/kjanat/chatlogger-api-go/test/fixtures"
	"github.com/kjanat/chatlogger-api-go/test/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewUserRepository(dbWrapper)

	// Create test organization first
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(t, err)

	// Create test user
	user := fixtures.CreateTestUser(org.ID)

	err = repo.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, org.ID, user.OrganizationID)
}

func TestUserRepository_FindByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewUserRepository(dbWrapper)

	// Create test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(t, err)

	user := fixtures.CreateTestUser(org.ID)
	err = repo.Create(user)
	require.NoError(t, err)

	// Test finding existing user
	foundUser, err := repo.FindByID(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, user.Role, foundUser.Role)

	// Test finding non-existent user
	nonExistentUser, err := repo.FindByID(9999)
	assert.NoError(t, err)
	assert.Nil(t, nonExistentUser)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewUserRepository(dbWrapper)

	// Create test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(t, err)

	user := fixtures.CreateTestUser(org.ID)
	user.Email = "unique@example.com"
	err = repo.Create(user)
	require.NoError(t, err)

	// Test finding existing user by email
	foundUser, err := repo.FindByEmail("unique@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, "unique@example.com", foundUser.Email)

	// Test finding non-existent user by email
	nonExistentUser, err := repo.FindByEmail("nonexistent@example.com")
	assert.NoError(t, err)
	assert.Nil(t, nonExistentUser)
}

func TestUserRepository_FindByOrganizationID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewUserRepository(dbWrapper)

	// Create test organizations
	org1 := fixtures.CreateTestOrganization()
	org1.Name = "Organization 1"
	err := db.Create(org1).Error
	require.NoError(t, err)

	org2 := fixtures.CreateTestOrganization()
	org2.Name = "Organization 2"
	err = db.Create(org2).Error
	require.NoError(t, err)

	// Create users for org1
	user1 := fixtures.CreateTestUser(org1.ID)
	user1.Email = "user1@org1.com"
	err = repo.Create(user1)
	require.NoError(t, err)

	user2 := fixtures.CreateTestUser(org1.ID)
	user2.Email = "user2@org1.com"
	err = repo.Create(user2)
	require.NoError(t, err)

	// Create user for org2
	user3 := fixtures.CreateTestUser(org2.ID)
	user3.Email = "user1@org2.com"
	err = repo.Create(user3)
	require.NoError(t, err)

	// Test finding users by organization ID
	org1Users, err := repo.FindByOrganizationID(org1.ID, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, org1Users, 2)

	org2Users, err := repo.FindByOrganizationID(org2.ID, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, org2Users, 1)

	// Test pagination
	org1UsersPage1, err := repo.FindByOrganizationID(org1.ID, 1, 0)
	assert.NoError(t, err)
	assert.Len(t, org1UsersPage1, 1)

	org1UsersPage2, err := repo.FindByOrganizationID(org1.ID, 1, 1)
	assert.NoError(t, err)
	assert.Len(t, org1UsersPage2, 1)
}

func TestUserRepository_Update(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewUserRepository(dbWrapper)

	// Create test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(t, err)

	user := fixtures.CreateTestUser(org.ID)
	err = repo.Create(user)
	require.NoError(t, err)

	// Update user
	user.FirstName = "Updated First"
	user.LastName = "Updated Last"
	user.Role = domain.RoleAdmin
	err = repo.Update(user)
	assert.NoError(t, err)

	// Verify update
	updatedUser, err := repo.FindByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated First", updatedUser.FirstName)
	assert.Equal(t, "Updated Last", updatedUser.LastName)
	assert.Equal(t, domain.RoleAdmin, updatedUser.Role)
}

func TestUserRepository_Delete(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewUserRepository(dbWrapper)

	// Create test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(t, err)

	user := fixtures.CreateTestUser(org.ID)
	err = repo.Create(user)
	require.NoError(t, err)

	// Delete user
	err = repo.Delete(user.ID)
	assert.NoError(t, err)

	// Verify deletion
	deletedUser, err := repo.FindByID(user.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedUser)
}

func TestUserRepository_EmailUniqueness(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)

	dbWrapper := &Database{DB: db}
	repo := NewUserRepository(dbWrapper)

	// Create test organization
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(t, err)

	// Create first user
	user1 := fixtures.CreateTestUser(org.ID)
	user1.Email = "duplicate@example.com"
	err = repo.Create(user1)
	require.NoError(t, err)

	// Try to create second user with same email
	user2 := fixtures.CreateTestUser(org.ID)
	user2.Email = "duplicate@example.com"
	err = repo.Create(user2)
	assert.Error(t, err) // Should fail due to unique constraint
}

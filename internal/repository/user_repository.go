package repository

import (
	"errors"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"gorm.io/gorm"
)

// UserRepo implements the domain.UserRepository interface.
type UserRepo struct {
	db *Database
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db *Database) domain.UserRepository {
	return &UserRepo{db: db}
}

// Create creates a new user.
func (r *UserRepo) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// FindByID finds a user by ID.
func (r *UserRepo) FindByID(id uint64) (*domain.User, error) {
	var user domain.User

	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

// FindByEmail finds a user by email.
func (r *UserRepo) FindByEmail(email string) (*domain.User, error) {
	var user domain.User

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

// FindByOrganizationID finds users by organization ID with pagination.
func (r *UserRepo) FindByOrganizationID(orgID uint64, limit, offset int) ([]domain.User, error) {
	var users []domain.User
	err := r.db.Where("organization_id = ?", orgID).Limit(limit).Offset(offset).Find(&users).Error

	return users, err
}

// Update updates a user.
func (r *UserRepo) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user by ID.
func (r *UserRepo) Delete(id uint64) error {
	return r.db.Delete(&domain.User{}, id).Error
}

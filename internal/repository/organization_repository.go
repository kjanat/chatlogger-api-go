package repository

import (
	"errors"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"gorm.io/gorm"
)

// OrganizationRepo implements the domain.OrganizationRepository interface.
type OrganizationRepo struct {
	db *Database
}

// NewOrganizationRepository creates a new organization repository.
func NewOrganizationRepository(db *Database) domain.OrganizationRepository {
	return &OrganizationRepo{db: db}
}

// Create creates a new organization.
func (r *OrganizationRepo) Create(org *domain.Organization) error {
	return r.db.Create(org).Error
}

// FindByID finds an organization by ID.
func (r *OrganizationRepo) FindByID(id uint64) (*domain.Organization, error) {
	var org domain.Organization

	err := r.db.First(&org, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &org, nil
}

// FindBySlug finds an organization by slug.
func (r *OrganizationRepo) FindBySlug(slug string) (*domain.Organization, error) {
	var org domain.Organization

	err := r.db.Where("slug = ?", slug).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &org, nil
}

// Update updates an organization.
func (r *OrganizationRepo) Update(org *domain.Organization) error {
	return r.db.Save(org).Error
}

// Delete deletes an organization by ID.
func (r *OrganizationRepo) Delete(id uint64) error {
	return r.db.Delete(&domain.Organization{}, id).Error
}

// List lists organizations with pagination.
func (r *OrganizationRepo) List(limit, offset int) ([]domain.Organization, error) {
	var orgs []domain.Organization
	err := r.db.Limit(limit).Offset(offset).Find(&orgs).Error

	return orgs, err
}

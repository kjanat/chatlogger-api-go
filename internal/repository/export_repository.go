package repository

import (
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"gorm.io/gorm"
)

// ExportRepository implements domain.ExportRepository for database operations.
type ExportRepository struct {
	db *gorm.DB
}

// NewExportRepository creates a new export repository.
func NewExportRepository(db *gorm.DB) *ExportRepository {
	return &ExportRepository{db: db}
}

// Create adds a new export record to the database.
func (r *ExportRepository) Create(export *domain.Export) error {
	return r.db.Create(export).Error
}

// GetByID retrieves an export by its ID.
func (r *ExportRepository) GetByID(id uint64) (*domain.Export, error) {
	var export domain.Export
	if err := r.db.First(&export, id).Error; err != nil {
		return nil, err
	}
	return &export, nil
}

// GetByOrganizationID retrieves exports for an organization.
func (r *ExportRepository) GetByOrganizationID(
	organizationID uint64,
	limit, offset int,
) ([]*domain.Export, error) {
	var exports []*domain.Export
	err := r.db.Where("organization_id = ?", organizationID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&exports).Error
	return exports, err
}

// UpdateStatus updates the status of an export.
func (r *ExportRepository) UpdateStatus(
	id uint64,
	status domain.ExportStatus,
	errorMsg string,
) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if status == domain.ExportStatusCompleted {
		now := time.Now()
		updates["completed_at"] = now
	}

	if errorMsg != "" {
		updates["error"] = errorMsg
	}

	return r.db.Model(&domain.Export{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateFilePath updates the file path of a completed export.
func (r *ExportRepository) UpdateFilePath(id uint64, filePath string) error {
	return r.db.Model(&domain.Export{}).Where("id = ?", id).Update("file_path", filePath).Error
}

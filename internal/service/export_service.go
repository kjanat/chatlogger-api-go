package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/kjanat/chatlogger-api-go/internal/jobs"
)

// ExportService implements domain.ExportService.
type ExportService struct {
	exportRepo domain.ExportRepository
	queue      *jobs.Queue
}

// NewExportService creates a new export service.
func NewExportService(exportRepo domain.ExportRepository, queue *jobs.Queue) *ExportService {
	return &ExportService{
		exportRepo: exportRepo,
		queue:      queue,
	}
}

// CreateExport creates an export job and enqueues it for processing.
func (s *ExportService) CreateExport(
	orgID, userID uint64,
	format domain.ExportFormat,
	exportType domain.ExportType,
) (*domain.Export, error) {
	// Create new export record
	export := &domain.Export{
		OrganizationID: orgID,
		UserID:         userID,
		Format:         format,
		Type:           exportType,
		Status:         domain.ExportStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Save to database
	if err := s.exportRepo.Create(export); err != nil {
		return nil, fmt.Errorf("failed to create export: %w", err)
	}

	// Enqueue the export job
	if err := s.queue.EnqueueExport(export.ID); err != nil {
		// If enqueueing fails, update the export status
		if updateErr := s.exportRepo.UpdateStatus(export.ID, domain.ExportStatusFailed, err.Error()); updateErr != nil {
			// Log the update error but return the original error
			return nil, fmt.Errorf(
				"failed to enqueue export job: %w (additionally, failed to update status: %v)",
				err,
				updateErr,
			)
		}
		return nil, fmt.Errorf("failed to enqueue export job: %w", err)
	}

	return export, nil
}

// GetExport gets an export by ID, ensuring it belongs to the given organization.
func (s *ExportService) GetExport(id, orgID uint64) (*domain.Export, error) {
	export, err := s.exportRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get export by ID: %w", err)
	}

	// Security check to ensure the export belongs to the organization
	if export.OrganizationID != orgID {
		return nil, errors.New("export not found")
	}

	return export, nil
}

// ListExports lists exports for an organization.
func (s *ExportService) ListExports(orgID uint64, limit, offset int) ([]*domain.Export, error) {
	exports, err := s.exportRepo.GetByOrganizationID(orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list exports for organization: %w", err)
	}
	return exports, nil
}

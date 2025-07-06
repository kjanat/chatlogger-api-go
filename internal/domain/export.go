package domain

import "time"

// ExportStatus represents the current status of an export job.
type ExportStatus string

const (
	ExportStatusPending    ExportStatus = "pending"
	ExportStatusProcessing ExportStatus = "processing"
	ExportStatusCompleted  ExportStatus = "completed"
	ExportStatusFailed     ExportStatus = "failed"
)

// ExportFormat represents the format of an export.
type ExportFormat string

const (
	ExportFormatJSON ExportFormat = "json"
	ExportFormatCSV  ExportFormat = "csv"
)

// ExportType represents the type of data being exported.
type ExportType string

const (
	ExportTypeChats    ExportType = "chats"
	ExportTypeMessages ExportType = "messages"
	ExportTypeAll      ExportType = "all"
)

// Export represents an asynchronous export job.
type Export struct {
	ID             uint64       `json:"id"                     gorm:"primaryKey"`
	OrganizationID uint64       `json:"organization_id"`
	UserID         uint64       `json:"user_id"`
	Format         ExportFormat `json:"format"`
	Type           ExportType   `json:"type"`
	Status         ExportStatus `json:"status"`
	FilePath       string       `json:"file_path,omitempty"`
	Error          string       `json:"error,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	CompletedAt    *time.Time   `json:"completed_at,omitempty"`
}

// ExportRepository defines the operations available on exports.
type ExportRepository interface {
	Create(export *Export) error
	GetByID(id uint64) (*Export, error)
	GetByOrganizationID(organizationID uint64, limit, offset int) ([]*Export, error)
	UpdateStatus(id uint64, status ExportStatus, errorMsg string) error
	UpdateFilePath(id uint64, filePath string) error
}

// ExportService defines the interface for export-related business logic.
type ExportService interface {
	CreateExport(orgID, userID uint64, format ExportFormat, exportType ExportType) (*Export, error)
	GetExport(id, orgID uint64) (*Export, error)
	ListExports(orgID uint64, limit, offset int) ([]*Export, error)
}

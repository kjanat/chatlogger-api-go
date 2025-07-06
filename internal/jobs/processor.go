// Package jobs provides asynchronous job processing capabilities for the ChatLogger API.
// It defines job types, payloads, and processors for handling background tasks
// such as exporting chat data in different formats (JSON, CSV).
// This file contains the export processor implementation which handles
// processing of export jobs asynchronously.
package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hibiken/asynq"
	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/kjanat/chatlogger-api-go/internal/strategy"
)

// ExportProcessor handles processing of export jobs.
type ExportProcessor struct {
	exportRepo     domain.ExportRepository
	chatService    domain.ChatService
	messageService domain.MessageService
	exportDir      string
}

// NewExportProcessor creates a new export processor.
func NewExportProcessor(
	exportRepo domain.ExportRepository,
	chatService domain.ChatService,
	messageService domain.MessageService,
	exportDir string,
) *ExportProcessor {
	return &ExportProcessor{
		exportRepo:     exportRepo,
		chatService:    chatService,
		messageService: messageService,
		exportDir:      exportDir,
	}
}

// ProcessExport processes an export job.
func (p *ExportProcessor) ProcessExport(ctx context.Context, task *asynq.Task) error {
	var payload ExportPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	export, err := p.exportRepo.GetByID(payload.ExportID)
	if err != nil {
		return fmt.Errorf("failed to get export: %w", err)
	}

	// Update status to processing
	if err := p.exportRepo.UpdateStatus(export.ID, domain.ExportStatusProcessing, ""); err != nil {
		return fmt.Errorf("failed to update export status: %w", err)
	}

	// Get chats for the organization
	chats, err := p.chatService.GetByOrganizationID(export.OrganizationID, 1000, 0)
	if err != nil {
		errorMsg := fmt.Sprintf("failed to get chats: %v", err)
		if err := p.exportRepo.UpdateStatus(export.ID, domain.ExportStatusFailed, errorMsg); err != nil {
			return fmt.Errorf("failed to update export status after chat error: %w", err)
		}
		return fmt.Errorf("%s", errorMsg)
	}

	// Load messages for each chat if needed
	if export.Type == domain.ExportTypeAll || export.Type == domain.ExportTypeMessages {
		for i := range chats {
			messages, err := p.messageService.GetByChatID(chats[i].ID)
			if err != nil {
				errorMsg := fmt.Sprintf("failed to get messages for chat %d: %v", chats[i].ID, err)
				if err := p.exportRepo.UpdateStatus(export.ID, domain.ExportStatusFailed, errorMsg); err != nil {
					return fmt.Errorf("failed to update export status after message error: %w", err)
				}
				return fmt.Errorf("%s", errorMsg)
			}
			chats[i].Messages = messages
		}
	}

	// Prepare data for export
	data := map[string]interface{}{
		"organization_id": export.OrganizationID,
		"export_date":     time.Now().Format(time.RFC3339),
		"export_type":     export.Type,
		"chats":           chats,
	}

	// Select appropriate exporter
	var exporter strategy.Exporter
	switch export.Format {
	case domain.ExportFormatJSON:
		exporter = &strategy.JSONExporter{}
	case domain.ExportFormatCSV:
		exporter = &strategy.CSVExporter{}
	default:
		errorMsg := "unsupported export format"
		if err := p.exportRepo.UpdateStatus(export.ID, domain.ExportStatusFailed, errorMsg); err != nil {
			return fmt.Errorf("failed to update export status for unsupported format: %w", err)
		}
		return fmt.Errorf("%s", errorMsg)
	}

	// Export the data
	exportData, err := exporter.Export(data)
	if err != nil {
		errorMsg := fmt.Sprintf("failed to export data: %v", err)
		if err := p.exportRepo.UpdateStatus(export.ID, domain.ExportStatusFailed, errorMsg); err != nil {
			return fmt.Errorf("failed to update export status after export error: %w", err)
		}
		return fmt.Errorf("%s", errorMsg)
	}

	// Create export directory if it doesn't exist
	if err := os.MkdirAll(p.exportDir, 0o750); err != nil {
		errorMsg := fmt.Sprintf("failed to create export directory: %v", err)
		if err := p.exportRepo.UpdateStatus(export.ID, domain.ExportStatusFailed, errorMsg); err != nil {
			return fmt.Errorf("failed to update export status after directory error: %w", err)
		}
		return fmt.Errorf("%s", errorMsg)
	}

	// Generate filename
	extension := ".json"
	if export.Format == domain.ExportFormatCSV {
		extension = ".csv"
	}

	filename := fmt.Sprintf("export_%d_%s%s",
		export.OrganizationID,
		time.Now().Format("20060102_150405"),
		extension)

	filePath := filepath.Join(p.exportDir, filename)

	// Write file
	if err := os.WriteFile(filePath, exportData, 0o600); err != nil {
		errorMsg := fmt.Sprintf("failed to write export file: %v", err)
		if err := p.exportRepo.UpdateStatus(export.ID, domain.ExportStatusFailed, errorMsg); err != nil {
			return fmt.Errorf("failed to update export status after file write error: %w", err)
		}
		return fmt.Errorf("%s", errorMsg)
	}

	// Update export record with file path
	if err := p.exportRepo.UpdateFilePath(export.ID, filePath); err != nil {
		return fmt.Errorf("failed to update file path: %w", err)
	}

	// Update status to completed
	if err := p.exportRepo.UpdateStatus(export.ID, domain.ExportStatusCompleted, ""); err != nil {
		return fmt.Errorf("failed to update export status to completed: %w", err)
	}
	return nil
}

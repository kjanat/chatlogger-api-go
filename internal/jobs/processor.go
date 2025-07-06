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
	export, err := p.parseAndInitializeExport(task)
	if err != nil {
		return err
	}

	chats, err := p.loadChatData(export)
	if err != nil {
		return err
	}

	filePath, err := p.generateExportFile(export, chats)
	if err != nil {
		return err
	}

	return p.finalizeExport(export, filePath)
}

// parseAndInitializeExport unmarshals the task payload and initializes the export.
func (p *ExportProcessor) parseAndInitializeExport(task *asynq.Task) (*domain.Export, error) {
	var payload ExportPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	export, err := p.exportRepo.GetByID(payload.ExportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get export: %w", err)
	}

	if err := p.exportRepo.UpdateStatus(export.ID, domain.ExportStatusProcessing, ""); err != nil {
		return nil, fmt.Errorf("failed to update export status: %w", err)
	}

	return export, nil
}

// loadChatData retrieves chats and messages for the export.
func (p *ExportProcessor) loadChatData(export *domain.Export) ([]domain.Chat, error) {
	chats, err := p.chatService.GetByOrganizationID(export.OrganizationID, 1000, 0)
	if err != nil {
		return nil, p.updateStatusAndReturnError(export.ID, "failed to get chats", err)
	}

	if export.Type == domain.ExportTypeAll || export.Type == domain.ExportTypeMessages {
		if err := p.loadMessagesForChats(export.ID, chats); err != nil {
			return nil, err
		}
	}

	return chats, nil
}

// loadMessagesForChats loads messages for each chat.
func (p *ExportProcessor) loadMessagesForChats(exportID uint64, chats []domain.Chat) error {
	for i := range chats {
		messages, err := p.messageService.GetByChatID(chats[i].ID)
		if err != nil {
			errorMsg := fmt.Sprintf("failed to get messages for chat %d", chats[i].ID)
			return p.updateStatusAndReturnError(exportID, errorMsg, err)
		}
		chats[i].Messages = messages
	}
	return nil
}

// generateExportFile creates the export file and returns the file path.
func (p *ExportProcessor) generateExportFile(
	export *domain.Export, chats []domain.Chat,
) (string, error) {
	data := map[string]interface{}{
		"organization_id": export.OrganizationID,
		"export_date":     time.Now().Format(time.RFC3339),
		"export_type":     export.Type,
		"chats":           chats,
	}

	exporter, err := p.createExporter(string(export.Format))
	if err != nil {
		return "", p.updateStatusAndReturnError(export.ID, "unsupported export format", err)
	}

	exportData, err := exporter.Export(data)
	if err != nil {
		return "", p.updateStatusAndReturnError(export.ID, "failed to export data", err)
	}

	return p.writeExportFile(export, exportData)
}

// createExporter returns the appropriate exporter for the format.
func (p *ExportProcessor) createExporter(format string) (strategy.Exporter, error) {
	switch format {
	case string(domain.ExportFormatJSON):
		return &strategy.JSONExporter{}, nil
	case string(domain.ExportFormatCSV):
		return &strategy.CSVExporter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// writeExportFile writes the export data to a file and returns the path.
func (p *ExportProcessor) writeExportFile(export *domain.Export, data []byte) (string, error) {
	if err := os.MkdirAll(p.exportDir, 0o750); err != nil {
		return "", p.updateStatusAndReturnError(export.ID, "failed to create export directory", err)
	}

	extension := ".json"
	if export.Format == domain.ExportFormatCSV {
		extension = ".csv"
	}

	filename := fmt.Sprintf("export_%d_%s%s",
		export.OrganizationID,
		time.Now().Format("20060102_150405"),
		extension)

	filePath := filepath.Join(p.exportDir, filename)

	if err := os.WriteFile(filePath, data, 0o600); err != nil {
		return "", p.updateStatusAndReturnError(export.ID, "failed to write export file", err)
	}

	return filePath, nil
}

// finalizeExport updates the export with the file path and marks it as completed.
func (p *ExportProcessor) finalizeExport(export *domain.Export, filePath string) error {
	if err := p.exportRepo.UpdateFilePath(export.ID, filePath); err != nil {
		return fmt.Errorf("failed to update file path: %w", err)
	}

	if err := p.exportRepo.UpdateStatus(export.ID, domain.ExportStatusCompleted, ""); err != nil {
		return fmt.Errorf("failed to update export status to completed: %w", err)
	}

	return nil
}

// updateStatusAndReturnError updates export status to failed and returns a formatted error.
func (p *ExportProcessor) updateStatusAndReturnError(
	exportID uint64, message string, err error,
) error {
	errorMsg := fmt.Sprintf("%s: %v", message, err)
	if updateErr := p.exportRepo.UpdateStatus(exportID, domain.ExportStatusFailed, errorMsg); updateErr != nil {
		return fmt.Errorf("failed to update export status after error (%s): %w", message, updateErr)
	}
	return fmt.Errorf("%s", errorMsg)
}

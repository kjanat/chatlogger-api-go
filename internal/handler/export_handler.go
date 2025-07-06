// Package handler provides HTTP request handlers for the ChatLogger API.
// This file implements handlers for chat data export functionality, supporting both
// synchronous and asynchronous export operations in different formats (JSON, CSV).
package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/kjanat/chatlogger-api-go/internal/middleware"
	"github.com/kjanat/chatlogger-api-go/internal/strategy"
)

// ExportHandler handles export-related requests.
type ExportHandler struct {
	exportService  domain.ExportService
	chatService    domain.ChatService
	messageService domain.MessageService
	exportDir      string
}

// NewExportHandler creates a new export handler.
func NewExportHandler(
	exportService domain.ExportService,
	chatService domain.ChatService,
	messageService domain.MessageService,
	exportDir string,
) *ExportHandler {
	return &ExportHandler{
		exportService:  exportService,
		chatService:    chatService,
		messageService: messageService,
		exportDir:      exportDir,
	}
}

// ExportRequest represents the request to export data.
type ExportRequest struct {
	Format string `binding:"required,oneof=json csv"           json:"format"`
	Type   string `binding:"required,oneof=chats messages all" json:"type"`
}

// CreateExport handles the request to create an asynchronous export
//
//	@Summary		Create Asynchronous Export Job
//	@Description	Initiates an asynchronous job to export chat data (chats, messages, or all) in JSON or CSV format for the user's organization.
//	@Tags			Exports
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ExportRequest			true	"Export parameters (format, type)"
//	@Success		202		{object}	map[string]interface{}	"export_id: uint64, status: domain.ExportStatus, message: string"
//	@Failure		400		{object}	map[string]string		"Invalid request data (format/type)"
//	@Failure		401		{object}	map[string]string		"Unauthorized (JWT invalid/missing, Org or User ID not found)"
//	@Failure		500		{object}	map[string]string		"Failed to create export job"
//	@Security		BearerAuth
//	@Router			/v1/exports [post]
func (h *ExportHandler) CreateExport(c *gin.Context) {
	var req ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Get organization ID from context
	orgID, exists := c.Get(middleware.OrganizationIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization not found"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Convert string format and type to domain types
	var format domain.ExportFormat
	switch req.Format {
	case "json":
		format = domain.ExportFormatJSON
	case "csv":
		format = domain.ExportFormatCSV
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported export format"})
		return
	}

	var exportType domain.ExportType
	switch req.Type {
	case "chats":
		exportType = domain.ExportTypeChats
	case "messages":
		exportType = domain.ExportTypeMessages
	case "all":
		exportType = domain.ExportTypeAll
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported export type"})
		return
	}

	// Type assert organization ID
	orgIDValue, ok := orgID.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid organization ID type"})
		return
	}

	// Type assert user ID
	userIDValue, ok := userID.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Create the export job
	export, err := h.exportService.CreateExport(
		orgIDValue,
		userIDValue,
		format,
		exportType,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to create export: " + err.Error()},
		)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"export_id": export.ID,
		"status":    export.Status,
		"message":   "Export job created successfully and queued for processing",
	})
}

// GetExport handles the request to get export status
//
//	@Summary		Get Export Job Status
//	@Description	Retrieves the status and details of a specific asynchronous export job.
//	@Tags			Exports
//	@Produce		json
//	@Param			id	path		uint64				true	"Export Job ID"
//	@Success		200	{object}	domain.Export		"Export job details"
//	@Failure		400	{object}	map[string]string	"Invalid export ID"
//	@Failure		401	{object}	map[string]string	"Unauthorized (JWT invalid/missing or Org ID not found)"
//	@Failure		404	{object}	map[string]string	"Export not found or doesn't belong to user's org"
//	@Security		BearerAuth
//	@Router			/v1/exports/{id} [get]
func (h *ExportHandler) GetExport(c *gin.Context) {
	exportID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid export ID"})
		return
	}

	// Get organization ID from context
	orgID, exists := c.Get(middleware.OrganizationIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization not found"})
		return
	}

	// Type assert organization ID
	orgIDValue, ok := orgID.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid organization ID type"})
		return
	}

	export, err := h.exportService.GetExport(exportID, orgIDValue)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Export not found"})
		return
	}

	c.JSON(http.StatusOK, export)
}

// DownloadExport handles the request to download a completed export
//
//	@Summary		Download Export File
//	@Description	Downloads the file generated by a completed asynchronous export job.
//	@Tags			Exports
//	@Produce		octet-stream
//	@Param			id	path		uint64				true	"Export Job ID"
//	@Success		200	{file}		file				"Export file (JSON or CSV)"
//	@Failure		400	{object}	map[string]string	"Invalid export ID or export not ready"
//	@Failure		401	{object}	map[string]string	"Unauthorized (JWT invalid/missing or Org ID not found)"
//	@Failure		404	{object}	map[string]string	"Export not found or doesn't belong to user's org"
//	@Failure		500	{object}	map[string]string	"Export file path not found or file missing on disk"
//	@Security		BearerAuth
//	@Router			/v1/exports/{id}/download [get]
func (h *ExportHandler) DownloadExport(c *gin.Context) {
	exportID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid export ID"})
		return
	}

	// Get organization ID from context
	orgID, exists := c.Get(middleware.OrganizationIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization not found"})
		return
	}

	// Type assert organization ID
	orgIDValue, ok := orgID.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid organization ID type"})
		return
	}

	export, err := h.exportService.GetExport(exportID, orgIDValue)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Export not found"})
		return
	}

	// Check if export is completed
	if export.Status != domain.ExportStatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Export is not ready for download",
			"status": export.Status,
		})
		return
	}

	// Check if file exists
	if export.FilePath == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Export file path not found"})
		return
	}

	// Check if file exists on disk
	if _, err := os.Stat(export.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Export file not found on disk"})
		return
	}

	// Set appropriate content type
	contentType := "application/json"
	if export.Format == domain.ExportFormatCSV {
		contentType = "text/csv"
	}

	filename := filepath.Base(export.FilePath)

	// Set headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", contentType)

	c.File(export.FilePath)
}

// ListExports handles the request to list exports for an organization
//
//	@Summary		List Export Jobs
//	@Description	Retrieves a paginated list of asynchronous export jobs for the user's organization.
//	@Tags			Exports
//	@Produce		json
//	@Param			limit	query		int					false	"Number of exports per page"	default(10)
//	@Param			offset	query		int					false	"Offset for pagination"			default(0)
//	@Success		200		{array}		domain.Export		"List of export jobs"
//	@Failure		401		{object}	map[string]string	"Unauthorized (JWT invalid/missing or Org ID not found)"
//	@Failure		500		{object}	map[string]string	"Failed to fetch exports"
//	@Security		BearerAuth
//	@Router			/v1/exports [get]
func (h *ExportHandler) ListExports(c *gin.Context) {
	// Get pagination parameters
	limit := 10
	offset := 0

	if limitParam := c.Query("limit"); limitParam != "" {
		if val, err := strconv.Atoi(limitParam); err == nil && val > 0 {
			limit = val
		}
	}

	if offsetParam := c.Query("offset"); offsetParam != "" {
		if val, err := strconv.Atoi(offsetParam); err == nil && val >= 0 {
			offset = val
		}
	}

	// Get organization ID from context
	orgID, exists := c.Get(middleware.OrganizationIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization not found"})
		return
	}

	// Type assert organization ID
	orgIDValue, ok := orgID.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid organization ID type"})
		return
	}

	exports, err := h.exportService.ListExports(orgIDValue, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exports"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exports": exports})
}

// SyncExport is the original synchronous export method, kept for backward compatibility
//
//	@Summary		Create Synchronous Export (Legacy)
//	@Description	Immediately generates and returns an export file (JSON or CSV) containing chat data. Use async export for large datasets.
//	@Tags			Exports (Legacy)
//	@Accept			json
//	@Produce		octet-stream
//	@Param			request	body		ExportRequest		true	"Export parameters (format, type)"
//	@Success		200		{file}		file				"Export file (JSON or CSV)"
//	@Failure		400		{object}	map[string]string	"Invalid request data (format/type)"
//	@Failure		401		{object}	map[string]string	"Unauthorized (JWT invalid/missing or Org ID not found)"
//	@Failure		500		{object}	map[string]string	"Failed to retrieve data or generate export"
//	@Security		BearerAuth
//	@Router			/v1/exports/sync [post]
func (h *ExportHandler) SyncExport(c *gin.Context) {
	req, orgID, err := h.parseSyncExportRequest(c)
	if err != nil {
		return // Error already sent in response
	}

	chats, err := h.loadSyncExportData(c, orgID, req)
	if err != nil {
		return // Error already sent in response
	}

	exportData, err := h.generateSyncExport(c, orgID, req, chats)
	if err != nil {
		return // Error already sent in response
	}

	h.sendSyncExportResponse(c, orgID, req, exportData)
}

// parseSyncExportRequest parses and validates the sync export request.
func (h *ExportHandler) parseSyncExportRequest(c *gin.Context) (*ExportRequest, uint64, error) {
	var req ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return nil, 0, fmt.Errorf("failed to bind JSON request: %w", err)
	}

	orgID, exists := c.Get(middleware.OrganizationIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return nil, 0, fmt.Errorf("missing org ID")
	}

	// Type assert organization ID
	orgIDValue, ok := orgID.(uint64)
	if !ok {
		return nil, 0, fmt.Errorf("invalid organization ID type")
	}

	return &req, orgIDValue, nil
}

// loadSyncExportData retrieves chats and messages for the export.
func (h *ExportHandler) loadSyncExportData(
	c *gin.Context, orgID uint64, req *ExportRequest,
) ([]domain.Chat, error) {
	chats, err := h.chatService.GetByOrganizationID(orgID, 1000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chats"})
		return nil, fmt.Errorf("failed to retrieve chats: %w", err)
	}

	if req.Type == "all" || req.Type == "messages" {
		if err := h.loadMessagesForSyncExport(c, chats); err != nil {
			return nil, err
		}
	}

	return chats, nil
}

// loadMessagesForSyncExport loads messages for each chat.
func (h *ExportHandler) loadMessagesForSyncExport(c *gin.Context, chats []domain.Chat) error {
	for i := range chats {
		messages, err := h.messageService.GetByChatID(chats[i].ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
			return fmt.Errorf("failed to retrieve messages for chat %d: %w", chats[i].ID, err)
		}
		chats[i].Messages = messages
	}
	return nil
}

// generateSyncExport creates the export data using the appropriate exporter.
func (h *ExportHandler) generateSyncExport(
	c *gin.Context,
	orgID uint64,
	req *ExportRequest,
	chats []domain.Chat,
) ([]byte, error) {
	data := gin.H{
		"organization_id": orgID,
		"export_date":     time.Now().Format(time.RFC3339),
		"export_type":     req.Type,
		"chats":           chats,
	}

	exporter, err := h.createSyncExporter(req.Format)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported export format"})
		return nil, err
	}

	exportData, err := exporter.Export(data)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to export data: " + err.Error()},
		)
		return nil, fmt.Errorf("failed to export data: %w", err)
	}

	return exportData, nil
}

// createSyncExporter returns the appropriate exporter for the format.
func (h *ExportHandler) createSyncExporter(format string) (strategy.Exporter, error) {
	switch format {
	case "json":
		return &strategy.JSONExporter{}, nil
	case "csv":
		return &strategy.CSVExporter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// sendSyncExportResponse sends the export data as a file download.
func (h *ExportHandler) sendSyncExportResponse(
	c *gin.Context,
	orgID uint64,
	req *ExportRequest,
	exportData []byte,
) {
	contentType, extension := h.getContentTypeAndExtension(req.Format)
	filename := fmt.Sprintf("chatlogger_export_%d_%s.%s",
		orgID, time.Now().Format("20060102150405"), extension)

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, contentType, exportData)
}

// getContentTypeAndExtension returns the appropriate content type and file extension for the format.
func (h *ExportHandler) getContentTypeAndExtension(format string) (string, string) {
	switch format {
	case "json":
		return "application/json", "json"
	case "csv":
		return "text/csv", "csv"
	default:
		return "application/octet-stream", "bin"
	}
}

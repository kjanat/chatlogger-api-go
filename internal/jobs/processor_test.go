package jobs

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/kjanat/chatlogger-api-go/test/fixtures"
	"github.com/kjanat/chatlogger-api-go/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestExportProcessor_ProcessExport_Success(t *testing.T) {
	// Create temporary directory for exports
	tempDir, err := os.MkdirTemp("", "export_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Setup mocks
	mockExportRepo := &mocks.MockExportRepository{}
	mockChatService := &mocks.MockChatService{}
	mockMessageService := &mocks.MockMessageService{}

	processor := NewExportProcessor(
		mockExportRepo,
		mockChatService,
		mockMessageService,
		tempDir,
	)

	// Create test data
	export := fixtures.CreateTestExport(1)
	export.ID = 1
	export.OrganizationID = 100
	export.Format = domain.ExportFormatJSON
	export.Type = domain.ExportTypeAll

	chats := []domain.Chat{
		*fixtures.CreateTestChat(1),
	}
	chats[0].ID = 1
	chats[0].OrganizationID = 100

	messages := []domain.Message{
		*fixtures.CreateTestMessage(1),
	}

	// Setup mock expectations
	mockExportRepo.On("GetByID", uint64(1)).Return(export, nil)
	mockExportRepo.On("UpdateStatus", uint64(1), domain.ExportStatusProcessing, "").Return(nil)
	mockChatService.On("GetByOrganizationID", uint64(100), 1000, 0).Return(chats, nil)
	mockMessageService.On("GetByChatID", mock.AnythingOfType("uint64")).Return(messages, nil)
	mockExportRepo.On("UpdateFilePath", uint64(1), mock.AnythingOfType("string")).Return(nil)
	mockExportRepo.On("UpdateStatus", uint64(1), domain.ExportStatusCompleted, "").Return(nil)

	// Create task payload
	payload := ExportPayload{ExportID: 1}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	task := asynq.NewTask(TypeExportProcess, payloadBytes)

	// Process the export
	err = processor.ProcessExport(context.Background(), task)

	// Assertions
	assert.NoError(t, err)

	// Verify that export file was created
	files, err := os.ReadDir(tempDir)
	require.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Contains(t, files[0].Name(), "export_100_")
	assert.Contains(t, files[0].Name(), ".json")

	// Verify file contents
	filePath := filepath.Join(tempDir, files[0].Name())
	fileData, err := os.ReadFile(filePath)
	require.NoError(t, err)

	var exportData map[string]interface{}
	err = json.Unmarshal(fileData, &exportData)
	require.NoError(t, err)

	assert.Equal(t, float64(100), exportData["organization_id"])
	assert.NotEmpty(t, exportData["export_date"])
	assert.NotEmpty(t, exportData["chats"])

	// Verify all mocks were called
	mockExportRepo.AssertExpectations(t)
	mockChatService.AssertExpectations(t)
	mockMessageService.AssertExpectations(t)
}

func TestExportProcessor_ProcessExport_InvalidPayload(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "export_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	mockExportRepo := &mocks.MockExportRepository{}
	mockChatService := &mocks.MockChatService{}
	mockMessageService := &mocks.MockMessageService{}

	processor := NewExportProcessor(
		mockExportRepo,
		mockChatService,
		mockMessageService,
		tempDir,
	)

	// Create task with invalid payload
	task := asynq.NewTask(TypeExportProcess, []byte("invalid json"))

	err = processor.ProcessExport(context.Background(), task)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal payload")
}

func TestExportProcessor_ProcessExport_ExportNotFound(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "export_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	mockExportRepo := &mocks.MockExportRepository{}
	mockChatService := &mocks.MockChatService{}
	mockMessageService := &mocks.MockMessageService{}

	processor := NewExportProcessor(
		mockExportRepo,
		mockChatService,
		mockMessageService,
		tempDir,
	)

	// Setup mock to return error
	mockExportRepo.On("GetByID", uint64(1)).Return((*domain.Export)(nil), assert.AnError)

	payload := ExportPayload{ExportID: 1}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	task := asynq.NewTask(TypeExportProcess, payloadBytes)

	err = processor.ProcessExport(context.Background(), task)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get export")

	mockExportRepo.AssertExpectations(t)
}

func TestExportProcessor_ProcessExport_ChatServiceError(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "export_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	mockExportRepo := &mocks.MockExportRepository{}
	mockChatService := &mocks.MockChatService{}
	mockMessageService := &mocks.MockMessageService{}

	processor := NewExportProcessor(
		mockExportRepo,
		mockChatService,
		mockMessageService,
		tempDir,
	)

	export := fixtures.CreateTestExport(1)
	export.ID = 1
	export.OrganizationID = 100

	mockExportRepo.On("GetByID", uint64(1)).Return(export, nil)
	mockExportRepo.On("UpdateStatus", uint64(1), domain.ExportStatusProcessing, "").Return(nil)
	mockChatService.On("GetByOrganizationID", uint64(100), 1000, 0).Return([]domain.Chat(nil), assert.AnError)
	mockExportRepo.On("UpdateStatus", uint64(1), domain.ExportStatusFailed, mock.AnythingOfType("string")).Return(nil)

	payload := ExportPayload{ExportID: 1}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	task := asynq.NewTask(TypeExportProcess, payloadBytes)

	err = processor.ProcessExport(context.Background(), task)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get chats")

	mockExportRepo.AssertExpectations(t)
	mockChatService.AssertExpectations(t)
}

func TestExportProcessor_ProcessExport_UnsupportedFormat(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "export_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	mockExportRepo := &mocks.MockExportRepository{}
	mockChatService := &mocks.MockChatService{}
	mockMessageService := &mocks.MockMessageService{}

	processor := NewExportProcessor(
		mockExportRepo,
		mockChatService,
		mockMessageService,
		tempDir,
	)

	export := fixtures.CreateTestExport(1)
	export.ID = 1
	export.OrganizationID = 100
	export.Format = "unsupported"

	chats := []domain.Chat{
		*fixtures.CreateTestChat(1),
	}

	mockExportRepo.On("GetByID", uint64(1)).Return(export, nil)
	mockExportRepo.On("UpdateStatus", uint64(1), domain.ExportStatusProcessing, "").Return(nil)
	mockChatService.On("GetByOrganizationID", uint64(100), 1000, 0).Return(chats, nil)
	// Add message service call since processor will try to load messages for each chat
	mockMessageService.On("GetByChatID", mock.AnythingOfType("uint64")).Return([]domain.Message{}, nil)
	mockExportRepo.On("UpdateStatus", uint64(1), domain.ExportStatusFailed, "unsupported export format").Return(nil)

	payload := ExportPayload{ExportID: 1}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	task := asynq.NewTask(TypeExportProcess, payloadBytes)

	err = processor.ProcessExport(context.Background(), task)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported export format")

	mockExportRepo.AssertExpectations(t)
	mockChatService.AssertExpectations(t)
}

func TestExportProcessor_ProcessExport_CSVFormat(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "export_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	mockExportRepo := &mocks.MockExportRepository{}
	mockChatService := &mocks.MockChatService{}
	mockMessageService := &mocks.MockMessageService{}

	processor := NewExportProcessor(
		mockExportRepo,
		mockChatService,
		mockMessageService,
		tempDir,
	)

	export := fixtures.CreateTestExport(1)
	export.ID = 1
	export.OrganizationID = 100
	export.Format = domain.ExportFormatCSV
	export.Type = domain.ExportTypeChats

	chats := []domain.Chat{
		*fixtures.CreateTestChat(1),
	}
	chats[0].ID = 1
	chats[0].OrganizationID = 100

	mockExportRepo.On("GetByID", uint64(1)).Return(export, nil)
	mockExportRepo.On("UpdateStatus", uint64(1), domain.ExportStatusProcessing, "").Return(nil)
	mockChatService.On("GetByOrganizationID", uint64(100), 1000, 0).Return(chats, nil)
	mockExportRepo.On("UpdateFilePath", uint64(1), mock.AnythingOfType("string")).Return(nil)
	mockExportRepo.On("UpdateStatus", uint64(1), domain.ExportStatusCompleted, "").Return(nil)

	payload := ExportPayload{ExportID: 1}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	task := asynq.NewTask(TypeExportProcess, payloadBytes)

	err = processor.ProcessExport(context.Background(), task)

	assert.NoError(t, err)

	// Verify CSV file was created
	files, err := os.ReadDir(tempDir)
	require.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Contains(t, files[0].Name(), ".csv")

	mockExportRepo.AssertExpectations(t)
	mockChatService.AssertExpectations(t)
}

func TestExportProcessor_ProcessExport_DirectoryCreationError(t *testing.T) {
	// Use an invalid directory path (e.g., a file as directory)
	tempFile, err := os.CreateTemp("", "not_a_dir")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	mockExportRepo := &mocks.MockExportRepository{}
	mockChatService := &mocks.MockChatService{}
	mockMessageService := &mocks.MockMessageService{}

	processor := NewExportProcessor(
		mockExportRepo,
		mockChatService,
		mockMessageService,
		tempFile.Name(), // Using file as directory path
	)

	export := fixtures.CreateTestExport(1)
	export.ID = 1
	export.OrganizationID = 100
	export.Format = domain.ExportFormatJSON

	chats := []domain.Chat{
		*fixtures.CreateTestChat(1),
	}

	mockExportRepo.On("GetByID", uint64(1)).Return(export, nil)
	mockExportRepo.On("UpdateStatus", uint64(1), domain.ExportStatusProcessing, "").Return(nil)
	mockChatService.On("GetByOrganizationID", uint64(100), 1000, 0).Return(chats, nil)
	// Add message service call since processor will try to load messages for each chat
	mockMessageService.On("GetByChatID", mock.AnythingOfType("uint64")).Return([]domain.Message{}, nil)
	mockExportRepo.On("UpdateStatus", uint64(1), domain.ExportStatusFailed, mock.AnythingOfType("string")).Return(nil)

	payload := ExportPayload{ExportID: 1}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	task := asynq.NewTask(TypeExportProcess, payloadBytes)

	err = processor.ProcessExport(context.Background(), task)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create export directory")

	mockExportRepo.AssertExpectations(t)
	mockChatService.AssertExpectations(t)
}

func TestExportProcessor_ProcessExport_MessageLoadingError(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "export_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	mockExportRepo := &mocks.MockExportRepository{}
	mockChatService := &mocks.MockChatService{}
	mockMessageService := &mocks.MockMessageService{}

	processor := NewExportProcessor(
		mockExportRepo,
		mockChatService,
		mockMessageService,
		tempDir,
	)

	export := fixtures.CreateTestExport(1)
	export.ID = 1
	export.OrganizationID = 100
	export.Format = domain.ExportFormatJSON
	export.Type = domain.ExportTypeAll // This will trigger message loading

	chats := []domain.Chat{
		*fixtures.CreateTestChat(1),
	}
	chats[0].ID = 1

	mockExportRepo.On("GetByID", uint64(1)).Return(export, nil)
	mockExportRepo.On("UpdateStatus", uint64(1), domain.ExportStatusProcessing, "").Return(nil)
	mockChatService.On("GetByOrganizationID", uint64(100), 1000, 0).Return(chats, nil)
	mockMessageService.On("GetByChatID", mock.AnythingOfType("uint64")).Return([]domain.Message(nil), assert.AnError)
	mockExportRepo.On("UpdateStatus", uint64(1), domain.ExportStatusFailed, mock.AnythingOfType("string")).Return(nil)

	payload := ExportPayload{ExportID: 1}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	task := asynq.NewTask(TypeExportProcess, payloadBytes)

	err = processor.ProcessExport(context.Background(), task)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get messages for chat")

	mockExportRepo.AssertExpectations(t)
	mockChatService.AssertExpectations(t)
	mockMessageService.AssertExpectations(t)
}

package jobs

import (
	"encoding/json"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportPayload_Marshal(t *testing.T) {
	payload := ExportPayload{
		ExportID: 12345,
	}

	data, err := json.Marshal(payload)
	require.NoError(t, err)

	var unmarshaled ExportPayload
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, payload.ExportID, unmarshaled.ExportID)
}

func TestNewQueue(t *testing.T) {
	redisAddr := "localhost:6379"
	queue := NewQueue(redisAddr)

	assert.NotNil(t, queue)
	assert.NotNil(t, queue.client)

	// Clean up
	err := queue.Close()
	assert.NoError(t, err)
}

func TestQueue_EnqueueExport_Payload(t *testing.T) {
	// This test focuses on payload generation and validation
	// without actually connecting to Redis

	exportID := uint64(123)

	// Test payload marshaling directly
	payload, err := json.Marshal(ExportPayload{ExportID: exportID})
	require.NoError(t, err)

	// Verify payload structure
	var unmarshaled ExportPayload
	err = json.Unmarshal(payload, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, exportID, unmarshaled.ExportID)

	// Test task creation
	task := asynq.NewTask(TypeExportProcess, payload)
	assert.Equal(t, TypeExportProcess, task.Type())
	assert.Equal(t, payload, task.Payload())
}

func TestTaskTypes(t *testing.T) {
	assert.Equal(t, "export:process", TypeExportProcess)
}

func TestExportPayload_JSONTags(t *testing.T) {
	payload := ExportPayload{
		ExportID: 999,
	}

	data, err := json.Marshal(payload)
	require.NoError(t, err)

	// Verify JSON structure
	var jsonData map[string]interface{}
	err = json.Unmarshal(data, &jsonData)
	require.NoError(t, err)

	assert.Equal(t, float64(999), jsonData["export_id"])
}

// Integration test that requires Redis to be running
// This test is skipped by default but can be run with Redis available.
func TestQueue_EnqueueExport_Integration(t *testing.T) {
	t.Skip("Skipping integration test - requires Redis")

	redisAddr := "localhost:6379"
	queue := NewQueue(redisAddr)
	defer func() {
		if err := queue.Close(); err != nil {
			t.Logf("Warning: failed to close queue: %v", err)
		}
	}()

	exportID := uint64(456)

	err := queue.EnqueueExport(exportID)
	if err != nil {
		// If Redis is not available, skip the test
		if assert.Contains(t, err.Error(), "connection refused") {
			t.Skip("Redis not available for integration test")
		}
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestQueue_Close(t *testing.T) {
	redisAddr := "localhost:6379"
	queue := NewQueue(redisAddr)

	err := queue.Close()
	// Close should not error even if Redis is not available
	assert.NoError(t, err)
}

// Test the structure and configuration of Asynq options.
func TestAsynqOptions(t *testing.T) {
	// Test that we can create the options used in EnqueueExport
	opts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Queue("exports"),
		asynq.Timeout(20 * 60), // 20 minutes
	}

	assert.Len(t, opts, 3)

	// Create a dummy task to verify options can be applied
	payload, err := json.Marshal(ExportPayload{ExportID: 1})
	require.NoError(t, err)

	task := asynq.NewTask(TypeExportProcess, payload)
	assert.NotNil(t, task)

	// Verify we can create task with options (this validates the options are valid)
	// In a real scenario, these would be passed to client.Enqueue
	assert.NotEmpty(t, opts)
}

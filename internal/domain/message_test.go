package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageRole_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		role     MessageRole
		expected bool
	}{
		{
			name:     "valid user role",
			role:     MessageRoleUser,
			expected: true,
		},
		{
			name:     "valid assistant role",
			role:     MessageRoleAssistant,
			expected: true,
		},
		{
			name:     "valid system role",
			role:     MessageRoleSystem,
			expected: true,
		},
		{
			name:     "invalid role",
			role:     MessageRole("invalid"),
			expected: false,
		},
		{
			name:     "empty role",
			role:     MessageRole(""),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.role.IsValid())
		})
	}
}

func TestMessageRole_Constants(t *testing.T) {
	assert.Equal(t, MessageRole("user"), MessageRoleUser)
	assert.Equal(t, MessageRole("assistant"), MessageRoleAssistant)
	assert.Equal(t, MessageRole("system"), MessageRoleSystem)
}

func TestMessage_GetMetadata(t *testing.T) {
	tests := []struct {
		name     string
		metadata string
		expected *MessageMetadata
		hasError bool
	}{
		{
			name:     "valid metadata",
			metadata: `{"token_count":150,"response_time":1250.5}`,
			expected: &MessageMetadata{
				TokenCount:   150,
				ResponseTime: 1250.5,
			},
			hasError: false,
		},
		{
			name:     "empty metadata",
			metadata: `{}`,
			expected: &MessageMetadata{},
			hasError: false,
		},
		{
			name:     "null metadata",
			metadata: "null",
			expected: &MessageMetadata{},
			hasError: false,
		},
		{
			name:     "empty string metadata",
			metadata: "",
			expected: &MessageMetadata{},
			hasError: false,
		},
		{
			name:     "invalid json metadata",
			metadata: `invalid json`,
			expected: nil,
			hasError: true,
		},
		{
			name:     "partial metadata",
			metadata: `{"token_count":75}`,
			expected: &MessageMetadata{
				TokenCount:   75,
				ResponseTime: 0,
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := &Message{Metadata: tt.metadata}
			result, err := message.GetMetadata()

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestMessage_SetMetadata(t *testing.T) {
	tests := []struct {
		name     string
		metadata *MessageMetadata
		expected string
		hasError bool
	}{
		{
			name: "valid metadata",
			metadata: &MessageMetadata{
				TokenCount:   150,
				ResponseTime: 1250.5,
			},
			expected: `{"token_count":150,"response_time":1250.5}`,
			hasError: false,
		},
		{
			name:     "empty metadata",
			metadata: &MessageMetadata{},
			expected: `{}`,
			hasError: false,
		},
		{
			name:     "nil metadata",
			metadata: nil,
			expected: `{}`,
			hasError: false,
		},
		{
			name: "partial metadata",
			metadata: &MessageMetadata{
				TokenCount: 100,
			},
			expected: `{"token_count":100}`,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := &Message{}
			err := message.SetMetadata(tt.metadata)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				
				// Parse both JSON strings to compare structure
				var expected, actual map[string]interface{}
				err1 := json.Unmarshal([]byte(tt.expected), &expected)
				err2 := json.Unmarshal([]byte(message.Metadata), &actual)
				
				require.NoError(t, err1)
				require.NoError(t, err2)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestMessage_Validate(t *testing.T) {
	tests := []struct {
		name     string
		message  *Message
		hasError bool
		errorMsg string
	}{
		{
			name: "valid message",
			message: &Message{
				Role:    MessageRoleUser,
				Content: "Hello, world!",
			},
			hasError: false,
		},
		{
			name: "invalid role",
			message: &Message{
				Role:    MessageRole("invalid"),
				Content: "Hello, world!",
			},
			hasError: true,
			errorMsg: "invalid message role, must be 'user', 'assistant', or 'system'",
		},
		{
			name: "empty content",
			message: &Message{
				Role:    MessageRoleUser,
				Content: "",
			},
			hasError: true,
			errorMsg: "message content cannot be empty",
		},
		{
			name: "valid assistant message",
			message: &Message{
				Role:    MessageRoleAssistant,
				Content: "I can help you with that.",
			},
			hasError: false,
		},
		{
			name: "valid system message",
			message: &Message{
				Role:    MessageRoleSystem,
				Content: "System initialization complete.",
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.Validate()

			if tt.hasError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMessage_MetadataRoundTrip(t *testing.T) {
	message := &Message{}
	originalMetadata := &MessageMetadata{
		TokenCount:   500,
		ResponseTime: 2500.75,
	}
	
	err := message.SetMetadata(originalMetadata)
	require.NoError(t, err)
	
	retrievedMetadata, err := message.GetMetadata()
	require.NoError(t, err)
	
	assert.Equal(t, originalMetadata, retrievedMetadata)
}

func TestMessage_Struct(t *testing.T) {
	now := time.Now()
	
	message := &Message{
		ID:        1,
		ChatID:    100,
		Role:      MessageRoleUser,
		Content:   "Test message content",
		Metadata:  `{"token_count":50}`,
		CreatedAt: now,
	}
	
	assert.Equal(t, uint64(1), message.ID)
	assert.Equal(t, uint64(100), message.ChatID)
	assert.Equal(t, MessageRoleUser, message.Role)
	assert.Equal(t, "Test message content", message.Content)
	assert.Equal(t, `{"token_count":50}`, message.Metadata)
	assert.Equal(t, now, message.CreatedAt)
}
package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChat_GetTags(t *testing.T) {
	tests := []struct {
		name     string
		tags     string
		expected []string
		hasError bool
	}{
		{
			name:     "valid tags",
			tags:     `["tag1", "tag2", "tag3"]`,
			expected: []string{"tag1", "tag2", "tag3"},
			hasError: false,
		},
		{
			name:     "empty tags",
			tags:     `[]`,
			expected: []string{},
			hasError: false,
		},
		{
			name:     "null tags",
			tags:     "null",
			expected: []string{},
			hasError: false,
		},
		{
			name:     "empty string tags",
			tags:     "",
			expected: []string{},
			hasError: false,
		},
		{
			name:     "invalid json",
			tags:     `invalid json`,
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chat := &Chat{Tags: tt.tags}
			result, err := chat.GetTags()

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestChat_SetTags(t *testing.T) {
	tests := []struct {
		name     string
		tags     []string
		expected string
		hasError bool
	}{
		{
			name:     "valid tags",
			tags:     []string{"tag1", "tag2", "tag3"},
			expected: `["tag1","tag2","tag3"]`,
			hasError: false,
		},
		{
			name:     "empty tags",
			tags:     []string{},
			expected: `[]`,
			hasError: false,
		},
		{
			name:     "nil tags",
			tags:     nil,
			expected: `[]`,
			hasError: false,
		},
		{
			name:     "single tag",
			tags:     []string{"solo"},
			expected: `["solo"]`,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chat := &Chat{}
			err := chat.SetTags(tt.tags)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, chat.Tags)
			}
		})
	}
}

func TestChat_GetMetadata(t *testing.T) {
	tests := []struct {
		name     string
		metadata string
		expected *ChatMetadata
		hasError bool
	}{
		{
			name:     "valid metadata",
			metadata: `{"ip_address":"192.168.1.1","country_code":"US","language_code":"en","sentiment":"positive","token_count":150}`,
			expected: &ChatMetadata{
				IPAddress:    "192.168.1.1",
				CountryCode:  "US",
				LanguageCode: "en",
				Sentiment:    "positive",
				TokenCount:   150,
			},
			hasError: false,
		},
		{
			name:     "empty metadata",
			metadata: `{}`,
			expected: &ChatMetadata{},
			hasError: false,
		},
		{
			name:     "null metadata",
			metadata: "null",
			expected: &ChatMetadata{},
			hasError: false,
		},
		{
			name:     "empty string metadata",
			metadata: "",
			expected: &ChatMetadata{},
			hasError: false,
		},
		{
			name:     "invalid json metadata",
			metadata: `invalid json`,
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chat := &Chat{Metadata: tt.metadata}
			result, err := chat.GetMetadata()

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestChat_SetMetadata(t *testing.T) {
	tests := []struct {
		name     string
		metadata *ChatMetadata
		expected string
		hasError bool
	}{
		{
			name: "valid metadata",
			metadata: &ChatMetadata{
				IPAddress:    "192.168.1.1",
				CountryCode:  "US",
				LanguageCode: "en",
				Sentiment:    "positive",
				TokenCount:   150,
			},
			expected: `{"ip_address":"192.168.1.1","country_code":"US","language_code":"en","sentiment":"positive","token_count":150}`,
			hasError: false,
		},
		{
			name:     "empty metadata",
			metadata: &ChatMetadata{},
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
			metadata: &ChatMetadata{
				IPAddress:   "10.0.0.1",
				Sentiment:   "neutral",
				TokenCount:  75,
				UserRating:  intPtr(4),
			},
			expected: `{"ip_address":"10.0.0.1","sentiment":"neutral","token_count":75,"user_rating":4}`,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chat := &Chat{}
			err := chat.SetMetadata(tt.metadata)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				
				// Parse both JSON strings to compare structure
				var expected, actual map[string]interface{}
				err1 := json.Unmarshal([]byte(tt.expected), &expected)
				err2 := json.Unmarshal([]byte(chat.Metadata), &actual)
				
				require.NoError(t, err1)
				require.NoError(t, err2)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestChat_TagsRoundTrip(t *testing.T) {
	chat := &Chat{}
	originalTags := []string{"test", "roundtrip", "tags"}
	
	err := chat.SetTags(originalTags)
	require.NoError(t, err)
	
	retrievedTags, err := chat.GetTags()
	require.NoError(t, err)
	
	assert.Equal(t, originalTags, retrievedTags)
}

func TestChat_MetadataRoundTrip(t *testing.T) {
	chat := &Chat{}
	originalMetadata := &ChatMetadata{
		IPAddress:        "192.168.1.100",
		CountryCode:      "CA",
		LanguageCode:     "fr",
		SessionID:        "session123",
		Sentiment:        "positive",
		IsEscalated:      true,
		IsForwardedToHR:  false,
		TranscriptLink:   "https://example.com/transcript",
		TokenCount:       300,
		AvgResponseTime:  2.5,
		QuestionCategory: "technical",
		UserRating:       intPtr(5),
	}
	
	err := chat.SetMetadata(originalMetadata)
	require.NoError(t, err)
	
	retrievedMetadata, err := chat.GetMetadata()
	require.NoError(t, err)
	
	assert.Equal(t, originalMetadata, retrievedMetadata)
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}
package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrganization_Struct(t *testing.T) {
	now := time.Now()

	org := &Organization{
		ID:        1,
		Name:      "Test Organization",
		Slug:      "test-org",
		Settings:  `{"feature_flags":{"export_enabled":true}}`,
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, uint64(1), org.ID)
	assert.Equal(t, "Test Organization", org.Name)
	assert.Equal(t, "test-org", org.Slug)
	assert.Equal(t, `{"feature_flags":{"export_enabled":true}}`, org.Settings)
	assert.Equal(t, now, org.CreatedAt)
	assert.Equal(t, now, org.UpdatedAt)
	assert.Empty(t, org.APIKeys)
	assert.Empty(t, org.Users)
	assert.Empty(t, org.Chats)
}

func TestOrganization_EmptySettings(t *testing.T) {
	org := &Organization{
		ID:        1,
		Name:      "Minimal Organization",
		Slug:      "minimal-org",
		Settings:  "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.Empty(t, org.Settings)
}

func TestOrganization_JSONSettings(t *testing.T) {
	testCases := []struct {
		name     string
		settings string
	}{
		{
			name:     "empty json object",
			settings: `{}`,
		},
		{
			name:     "complex settings",
			settings: `{"features":{"export":true,"api_rate_limit":1000},"branding":{"logo_url":"https://example.com/logo.png"}}`,
		},
		{
			name:     "null settings",
			settings: "null",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			org := &Organization{
				ID:        1,
				Name:      "Test Org",
				Slug:      "test-org",
				Settings:  tc.settings,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			assert.Equal(t, tc.settings, org.Settings)
		})
	}
}

func TestOrganization_SlugValidation(t *testing.T) {
	testCases := []struct {
		name string
		slug string
	}{
		{
			name: "valid slug with hyphen",
			slug: "my-organization",
		},
		{
			name: "valid slug lowercase",
			slug: "myorganization",
		},
		{
			name: "valid slug with numbers",
			slug: "org123",
		},
		{
			name: "short slug",
			slug: "org",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			org := &Organization{
				ID:        1,
				Name:      "Test Organization",
				Slug:      tc.slug,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			assert.Equal(t, tc.slug, org.Slug)
		})
	}
}

func TestOrganization_Relationships(t *testing.T) {
	org := &Organization{
		ID:        1,
		Name:      "Test Organization",
		Slug:      "test-org",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		APIKeys:   []APIKey{},
		Users:     []User{},
		Chats:     []Chat{},
	}

	// Test that relationship slices are initialized but empty
	assert.NotNil(t, org.APIKeys)
	assert.NotNil(t, org.Users)
	assert.NotNil(t, org.Chats)
	assert.Len(t, org.APIKeys, 0)
	assert.Len(t, org.Users, 0)
	assert.Len(t, org.Chats, 0)
}

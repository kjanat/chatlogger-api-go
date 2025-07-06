// Package service implements the business logic layer for the ChatLogger API,
// providing implementations of the domain service interfaces.
package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
)

// APIKeyService implements the domain.APIKeyService interface.
type APIKeyService struct {
	apiKeyRepo domain.APIKeyRepository
}

// NewAPIKeyService creates a new API key service.
func NewAPIKeyService(apiKeyRepo domain.APIKeyRepository) domain.APIKeyService {
	return &APIKeyService{
		apiKeyRepo: apiKeyRepo,
	}
}

// GenerateKey generates a new API key for an organization.
func (s *APIKeyService) GenerateKey(orgID uint64, label string) (string, error) {
	// Generate a random key
	rawBytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(rawBytes); err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}

	// Encode as base64 for easier API key usage
	rawKey := base64.URLEncoding.EncodeToString(rawBytes)

	// Hash the key for storage
	hashedKey := hashKey(rawKey)

	// Create the API key record
	key := &domain.APIKey{
		OrganizationID: orgID,
		HashedKey:      hashedKey,
		Label:          label,
		CreatedAt:      time.Now(),
	}

	if err := s.apiKeyRepo.Create(key); err != nil {
		return "", fmt.Errorf("failed to store API key: %w", err)
	}

	// Return the raw key - this is the only time it will be available
	return rawKey, nil
}

// ValidateKey validates a raw API key and returns the associated API key if valid.
func (s *APIKeyService) ValidateKey(rawKey string) (*domain.APIKey, error) {
	// Hash the key for lookup
	hashedKey := hashKey(rawKey)

	// Find the key by its hash
	key, err := s.apiKeyRepo.FindByHashedKey(hashedKey)
	if err != nil {
		return nil, fmt.Errorf("error looking up API key: %w", err)
	}

	if key == nil {
		return nil, nil // Key not found or revoked
	}

	return key, nil
}

// GetByID gets an API key by ID.
func (s *APIKeyService) GetByID(id uint64) (*domain.APIKey, error) {
	key, err := s.apiKeyRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find API key by ID: %w", err)
	}
	return key, nil
}

// ListByOrganizationID lists API keys for an organization.
func (s *APIKeyService) ListByOrganizationID(orgID uint64) ([]domain.APIKey, error) {
	keys, err := s.apiKeyRepo.ListByOrganizationID(orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys for organization: %w", err)
	}
	return keys, nil
}

// RevokeKey revokes an API key.
func (s *APIKeyService) RevokeKey(id uint64) error {
	if err := s.apiKeyRepo.Revoke(id); err != nil {
		return fmt.Errorf("failed to revoke API key: %w", err)
	}
	return nil
}

// DeleteKey permanently deletes an API key.
func (s *APIKeyService) DeleteKey(id uint64) error {
	if err := s.apiKeyRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}
	return nil
}

// hashKey hashes an API key for secure storage.
func hashKey(key string) string {
	h := sha256.New()
	h.Write([]byte(key))

	return hex.EncodeToString(h.Sum(nil))
}

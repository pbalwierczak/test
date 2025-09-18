package apikey

import (
	"crypto/subtle"
	"errors"
	"strings"
)

// Validator handles API key validation
type Validator struct {
	validKey string
}

// NewValidator creates a new API key validator
func NewValidator(apiKey string) *Validator {
	return &Validator{
		validKey: apiKey,
	}
}

// ValidateAPIKey validates the provided API key against the configured key
func (v *Validator) ValidateAPIKey(providedKey string) error {
	if providedKey == "" {
		return errors.New("API key is required")
	}

	// Use constant time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(providedKey), []byte(v.validKey)) != 1 {
		return errors.New("invalid API key")
	}

	return nil
}

// ExtractAPIKey extracts the API key from the Authorization header
// Supports both "Bearer <key>" and direct key formats
func ExtractAPIKey(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	// Check if it's a Bearer token
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return strings.TrimSpace(authHeader[7:])
	}

	// Return the header as-is if it's not a Bearer token
	return strings.TrimSpace(authHeader)
}

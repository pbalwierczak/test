package apikey

import (
	"crypto/subtle"
	"errors"
	"strings"
)

type Validator struct {
	validKey string
}

func NewValidator(apiKey string) *Validator {
	return &Validator{
		validKey: apiKey,
	}
}

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

func ExtractAPIKey(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return strings.TrimSpace(authHeader[7:])
	}

	return strings.TrimSpace(authHeader)
}

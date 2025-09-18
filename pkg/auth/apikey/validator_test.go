package apikey

import (
	"testing"
)

func TestValidator_ValidateAPIKey(t *testing.T) {
	validKey := "test-api-key-12345"
	validator := NewValidator(validKey)

	tests := []struct {
		name        string
		providedKey string
		expectError bool
	}{
		{
			name:        "valid API key",
			providedKey: "test-api-key-12345",
			expectError: false,
		},
		{
			name:        "invalid API key",
			providedKey: "wrong-key",
			expectError: true,
		},
		{
			name:        "empty API key",
			providedKey: "",
			expectError: true,
		},
		{
			name:        "case sensitive",
			providedKey: "TEST-API-KEY-12345",
			expectError: true,
		},
		{
			name:        "partial match",
			providedKey: "test-api-key-1234",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateAPIKey(tt.providedKey)
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestExtractAPIKey(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		expectedKey string
	}{
		{
			name:        "bearer token format",
			authHeader:  "Bearer test-api-key-12345",
			expectedKey: "test-api-key-12345",
		},
		{
			name:        "bearer token with spaces",
			authHeader:  "Bearer  test-api-key-12345  ",
			expectedKey: "test-api-key-12345",
		},
		{
			name:        "direct key format",
			authHeader:  "test-api-key-12345",
			expectedKey: "test-api-key-12345",
		},
		{
			name:        "direct key with spaces",
			authHeader:  "  test-api-key-12345  ",
			expectedKey: "test-api-key-12345",
		},
		{
			name:        "case insensitive bearer",
			authHeader:  "bearer test-api-key-12345",
			expectedKey: "test-api-key-12345",
		},
		{
			name:        "empty header",
			authHeader:  "",
			expectedKey: "",
		},
		{
			name:        "only bearer",
			authHeader:  "Bearer",
			expectedKey: "Bearer",
		},
		{
			name:        "bearer with empty key",
			authHeader:  "Bearer ",
			expectedKey: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractAPIKey(tt.authHeader)
			if result != tt.expectedKey {
				t.Errorf("Expected %q, got %q", tt.expectedKey, result)
			}
		})
	}
}

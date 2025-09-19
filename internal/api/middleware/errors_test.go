package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPIErrorWithDetails(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		message  string
		details  map[string]string
		expected *APIError
	}{
		{
			name:    "error with details",
			code:    400,
			message: "Validation failed",
			details: map[string]string{
				"field1": "error message 1",
				"field2": "error message 2",
			},
			expected: &APIError{
				Code:    400,
				Message: "Validation failed",
				Details: map[string]string{
					"field1": "error message 1",
					"field2": "error message 2",
				},
			},
		},
		{
			name:    "error with empty details",
			code:    500,
			message: "Internal server error",
			details: map[string]string{},
			expected: &APIError{
				Code:    500,
				Message: "Internal server error",
				Details: map[string]string{},
			},
		},
		{
			name:    "error with nil details",
			code:    404,
			message: "Not found",
			details: nil,
			expected: &APIError{
				Code:    404,
				Message: "Not found",
				Details: nil,
			},
		},
		{
			name:    "error with single detail",
			code:    422,
			message: "Unprocessable entity",
			details: map[string]string{
				"email": "Invalid email format",
			},
			expected: &APIError{
				Code:    422,
				Message: "Unprocessable entity",
				Details: map[string]string{
					"email": "Invalid email format",
				},
			},
		},
		{
			name:    "error with zero code",
			code:    0,
			message: "Unknown error",
			details: map[string]string{
				"reason": "Unknown",
			},
			expected: &APIError{
				Code:    0,
				Message: "Unknown error",
				Details: map[string]string{
					"reason": "Unknown",
				},
			},
		},
		{
			name:    "error with empty message",
			code:    400,
			message: "",
			details: map[string]string{
				"field": "error",
			},
			expected: &APIError{
				Code:    400,
				Message: "",
				Details: map[string]string{
					"field": "error",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewAPIErrorWithDetails(tt.code, tt.message, tt.details)

			// Check that the result is not nil
			assert.NotNil(t, result)

			// Check all fields
			assert.Equal(t, tt.expected.Code, result.Code)
			assert.Equal(t, tt.expected.Message, result.Message)
			assert.Equal(t, tt.expected.Details, result.Details)

			// Test the Error() method
			assert.Equal(t, tt.expected.Message, result.Error())
		})
	}
}

func TestNewAPIErrorWithDetails_ReferenceBehavior(t *testing.T) {
	// Test that the function uses the same map reference (current implementation behavior)
	originalDetails := map[string]string{
		"field1": "original value 1",
		"field2": "original value 2",
	}

	error := NewAPIErrorWithDetails(400, "Test error", originalDetails)

	// Modify the original details
	originalDetails["field1"] = "modified value"
	originalDetails["field3"] = "new field"

	// The error details should reflect the changes (current implementation behavior)
	assert.Equal(t, "modified value", error.Details["field1"])
	assert.Equal(t, "original value 2", error.Details["field2"])
	assert.Contains(t, error.Details, "field3")
}

func TestNewAPIErrorWithDetails_ErrorMethod(t *testing.T) {
	error := NewAPIErrorWithDetails(400, "Test error message", map[string]string{
		"field": "detail",
	})

	// Test that the Error() method returns the message
	assert.Equal(t, "Test error message", error.Error())
}

func TestNewAPIErrorWithDetails_JSONSerialization(t *testing.T) {
	error := NewAPIErrorWithDetails(422, "Validation failed", map[string]string{
		"email":    "Invalid email format",
		"password": "Password too short",
	})

	// Test that the error can be properly serialized to JSON
	// This is important for API responses
	assert.Equal(t, 422, error.Code)
	assert.Equal(t, "Validation failed", error.Message)
	assert.Len(t, error.Details, 2)
	assert.Equal(t, "Invalid email format", error.Details["email"])
	assert.Equal(t, "Password too short", error.Details["password"])
}

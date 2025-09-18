package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"scootin-aboot/pkg/auth/apikey"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAPIKeyMiddleware(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	validKey := "test-api-key-12345"
	validator := apikey.NewValidator(validKey)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "valid API key with Bearer",
			authHeader:     "Bearer test-api-key-12345",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "valid API key without Bearer",
			authHeader:     "test-api-key-12345",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "invalid API key",
			authHeader:     "Bearer wrong-key",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:           "missing API key",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:           "empty Bearer token",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test router
			router := gin.New()
			router.Use(APIKeyMiddleware(validator))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Create a test request
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Create a response recorder
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Check the response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectError {
				assert.Contains(t, w.Body.String(), "Authentication failed")
			} else {
				assert.Contains(t, w.Body.String(), "success")
			}
		})
	}
}

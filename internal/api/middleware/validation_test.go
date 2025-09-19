package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidateRequiredHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		requiredHeaders []string
		headers         map[string]string
		expectedStatus  int
		expectedError   string
	}{
		{
			name:            "all required headers present",
			requiredHeaders: []string{"Authorization", "Content-Type"},
			headers: map[string]string{
				"Authorization": "Bearer token123",
				"Content-Type":  "application/json",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:            "missing one required header",
			requiredHeaders: []string{"Authorization", "Content-Type"},
			headers: map[string]string{
				"Authorization": "Bearer token123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Missing Required Headers",
		},
		{
			name:            "missing multiple required headers",
			requiredHeaders: []string{"Authorization", "Content-Type", "X-API-Key"},
			headers: map[string]string{
				"Authorization": "Bearer token123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Missing Required Headers",
		},
		{
			name:            "empty required headers list",
			requiredHeaders: []string{},
			headers: map[string]string{
				"Authorization": "Bearer token123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:            "header with empty value",
			requiredHeaders: []string{"Authorization"},
			headers: map[string]string{
				"Authorization": "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Missing Required Headers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Gin router
			router := gin.New()

			// Add the validation middleware
			router.Use(ValidateRequiredHeaders(tt.requiredHeaders))

			// Add a test endpoint
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Create a test request
			req := httptest.NewRequest("GET", "/test", nil)

			// Add headers to the request
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			// Create a response recorder
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Assert the status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// If we expect an error, check the response body
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
				assert.Contains(t, w.Body.String(), "One or more required headers are missing")
			} else {
				assert.Contains(t, w.Body.String(), "success")
			}
		})
	}
}

func TestValidateContentLength(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		maxSize        int64
		contentLength  int64
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "content length within limit",
			maxSize:        1024,
			contentLength:  512,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "content length at limit",
			maxSize:        1024,
			contentLength:  1024,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "content length exceeds limit",
			maxSize:        1024,
			contentLength:  2048,
			expectedStatus: http.StatusRequestEntityTooLarge,
			expectedError:  "Request Too Large",
		},
		{
			name:           "zero content length",
			maxSize:        1024,
			contentLength:  0,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "very small limit",
			maxSize:        10,
			contentLength:  5,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "very small limit exceeded",
			maxSize:        10,
			contentLength:  15,
			expectedStatus: http.StatusRequestEntityTooLarge,
			expectedError:  "Request Too Large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Gin router
			router := gin.New()

			// Add the validation middleware
			router.Use(ValidateContentLength(tt.maxSize))

			// Add a test endpoint
			router.POST("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Create a test request
			req := httptest.NewRequest("POST", "/test", nil)
			req.ContentLength = tt.contentLength

			// Create a response recorder
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Assert the status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// If we expect an error, check the response body
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
				assert.Contains(t, w.Body.String(), "Request body exceeds maximum allowed size")
			} else {
				assert.Contains(t, w.Body.String(), "success")
			}
		})
	}
}

func TestValidateContentLength_GETRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test that GET requests are also affected by content length validation
	router := gin.New()
	router.Use(ValidateContentLength(10)) // Very small limit

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.ContentLength = 1000 // This should trigger the validation

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
	assert.Contains(t, w.Body.String(), "Request Too Large")
}

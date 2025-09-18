package middleware

import (
	"net/http"

	"scootin-aboot/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Code    int               `json:"code"`
	Details map[string]string `json:"details,omitempty"`
}

// ErrorHandlerMiddleware handles errors and provides consistent error responses
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle errors that were set by handlers
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Log the error
			utils.Error("Request error",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("client_ip", c.ClientIP()),
				zap.Error(err.Err),
			)

			// Determine the appropriate status code and response
			statusCode := http.StatusInternalServerError
			message := "Internal server error"

			// Check if it's a custom error type
			if customErr, ok := err.Err.(*APIError); ok {
				statusCode = customErr.Code
				message = customErr.Message
			} else {
				// Handle common error types
				switch err.Type {
				case gin.ErrorTypeBind:
					statusCode = http.StatusBadRequest
					message = "Invalid request parameters"
				case gin.ErrorTypeRender:
					statusCode = http.StatusInternalServerError
					message = "Error rendering response"
				case gin.ErrorTypePublic:
					statusCode = http.StatusBadRequest
					message = err.Error()
				}
			}

			// Send error response
			c.JSON(statusCode, ErrorResponse{
				Error:   http.StatusText(statusCode),
				Message: message,
				Code:    statusCode,
			})
		}
	}
}

// APIError represents a custom API error
type APIError struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new API error
func NewAPIError(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// NewAPIErrorWithDetails creates a new API error with additional details
func NewAPIErrorWithDetails(code int, message string, details map[string]string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Common error responses
var (
	ErrInvalidRequest      = NewAPIError(http.StatusBadRequest, "Invalid request")
	ErrUnauthorized        = NewAPIError(http.StatusUnauthorized, "Unauthorized")
	ErrForbidden           = NewAPIError(http.StatusForbidden, "Forbidden")
	ErrNotFound            = NewAPIError(http.StatusNotFound, "Resource not found")
	ErrMethodNotAllowed    = NewAPIError(http.StatusMethodNotAllowed, "Method not allowed")
	ErrConflict            = NewAPIError(http.StatusConflict, "Resource conflict")
	ErrUnprocessableEntity = NewAPIError(http.StatusUnprocessableEntity, "Unprocessable entity")
	ErrTooManyRequests     = NewAPIError(http.StatusTooManyRequests, "Too many requests")
	ErrInternalServer      = NewAPIError(http.StatusInternalServerError, "Internal server error")
	ErrServiceUnavailable  = NewAPIError(http.StatusServiceUnavailable, "Service unavailable")
)

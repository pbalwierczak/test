package middleware

import (
	"net/http"

	"scootin-aboot/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Code    int               `json:"code"`
	Details map[string]string `json:"details,omitempty"`
}

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			utils.Error("Request error",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("client_ip", c.ClientIP()),
				zap.Error(err.Err),
			)

			statusCode := http.StatusInternalServerError
			message := "Internal server error"

			if customErr, ok := err.Err.(*APIError); ok {
				statusCode = customErr.Code
				message = customErr.Message
			} else {
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

			c.JSON(statusCode, ErrorResponse{
				Error:   http.StatusText(statusCode),
				Message: message,
				Code:    statusCode,
			})
		}
	}
}

type APIError struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

func NewAPIErrorWithDetails(code int, message string, details map[string]string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

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

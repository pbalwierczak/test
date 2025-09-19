package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Code    int               `json:"code"`
	Details []ValidationError `json:"details"`
}

func ValidateJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if contentType != "" && contentType != "application/json" && contentType != "application/json; charset=utf-8" {
				c.JSON(http.StatusBadRequest, ValidationErrorResponse{
					Error:   "Invalid Content-Type",
					Message: "Content-Type must be application/json",
					Code:    http.StatusBadRequest,
					Details: []ValidationError{
						{
							Field:   "Content-Type",
							Message: "Must be application/json",
							Value:   contentType,
						},
					},
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func ValidateRequiredHeaders(requiredHeaders []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var missingHeaders []ValidationError

		for _, header := range requiredHeaders {
			if c.GetHeader(header) == "" {
				missingHeaders = append(missingHeaders, ValidationError{
					Field:   header,
					Message: "Required header is missing",
				})
			}
		}

		if len(missingHeaders) > 0 {
			c.JSON(http.StatusBadRequest, ValidationErrorResponse{
				Error:   "Missing Required Headers",
				Message: "One or more required headers are missing",
				Code:    http.StatusBadRequest,
				Details: missingHeaders,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func ValidateContentLength(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, ValidationErrorResponse{
				Error:   "Request Too Large",
				Message: "Request body exceeds maximum allowed size",
				Code:    http.StatusRequestEntityTooLarge,
				Details: []ValidationError{
					{
						Field:   "Content-Length",
						Message: "Request body too large",
						Value:   string(rune(c.Request.ContentLength)),
					},
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

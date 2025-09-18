package middleware

import (
	"net/http"

	"scootin-aboot/pkg/auth/apikey"

	"github.com/gin-gonic/gin"
)

// APIKeyMiddleware creates middleware for API key authentication
func APIKeyMiddleware(validator *apikey.Validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract API key from Authorization header
		authHeader := c.GetHeader("Authorization")
		apiKey := apikey.ExtractAPIKey(authHeader)

		// Validate the API key
		if err := validator.ValidateAPIKey(apiKey); err != nil {
			// Return unauthorized response directly
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Authentication failed",
				"code":    http.StatusUnauthorized,
				"details": map[string]string{
					"reason": err.Error(),
				},
			})
			c.Abort()
			return
		}

		// API key is valid, continue to next handler
		c.Next()
	}
}

// OptionalAPIKeyMiddleware creates middleware for optional API key authentication
// This can be used for endpoints that work with or without authentication
func OptionalAPIKeyMiddleware(validator *apikey.Validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract API key from Authorization header
		authHeader := c.GetHeader("Authorization")
		apiKey := apikey.ExtractAPIKey(authHeader)

		// If no Authorization header provided, continue without authentication
		if authHeader == "" {
			c.Next()
			return
		}

		// If Authorization header is provided but no valid API key, fail
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Authentication failed",
				"code":    http.StatusUnauthorized,
				"details": map[string]string{
					"reason": "API key is required",
				},
			})
			c.Abort()
			return
		}

		// Validate the API key if provided
		if err := validator.ValidateAPIKey(apiKey); err != nil {
			// Return unauthorized response directly
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Authentication failed",
				"code":    http.StatusUnauthorized,
				"details": map[string]string{
					"reason": err.Error(),
				},
			})
			c.Abort()
			return
		}

		// API key is valid, continue to next handler
		c.Next()
	}
}

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

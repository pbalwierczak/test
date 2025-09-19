package middleware

import (
	"net/http"

	"scootin-aboot/pkg/auth/apikey"

	"github.com/gin-gonic/gin"
)

func APIKeyMiddleware(validator *apikey.Validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var apiKey string

		if apiKeyHeader := c.GetHeader("X-API-Key"); apiKeyHeader != "" {
			apiKey = apiKeyHeader
		} else {
			authHeader := c.GetHeader("Authorization")
			apiKey = apikey.ExtractAPIKey(authHeader)
		}

		if err := validator.ValidateAPIKey(apiKey); err != nil {
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

		c.Next()
	}
}

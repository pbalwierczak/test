package routes

import (
	"scootin-aboot/internal/api/handlers"
	"scootin-aboot/internal/api/middleware"
	"scootin-aboot/pkg/auth/apikey"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, apiKey string) {
	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()

	// Initialize API key validator
	apiKeyValidator := apikey.NewValidator(apiKey)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Health check endpoint (no authentication required)
		v1.GET("/health", healthHandler.HealthCheck)

		// Protected routes (require API key authentication)
		protected := v1.Group("")
		protected.Use(middleware.APIKeyMiddleware(apiKeyValidator))
		{
			// Scooter management endpoints will be added here
			// protected.POST("/scooters/:id/trip/start", ...)
			// protected.POST("/scooters/:id/trip/end", ...)
			// protected.POST("/scooters/:id/location", ...)
			// protected.GET("/scooters", ...)
			// protected.GET("/scooters/:id", ...)
			// protected.GET("/scooters/closest", ...)
		}
	}
}

package routes

import (
	"scootin-aboot/internal/api/handlers"
	"scootin-aboot/internal/api/middleware"
	"scootin-aboot/internal/services"
	"scootin-aboot/pkg/auth/apikey"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, apiKey string, tripService services.TripService) {
	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	scooterHandler := handlers.NewScooterHandler(tripService)

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
			// Scooter management endpoints
			protected.POST("/scooters/:id/trip/start", scooterHandler.StartTrip)
			protected.POST("/scooters/:id/trip/end", scooterHandler.EndTrip)
			protected.POST("/scooters/:id/location", scooterHandler.UpdateLocation)
			protected.GET("/scooters", scooterHandler.GetScooters)
			protected.GET("/scooters/:id", scooterHandler.GetScooter)
			protected.GET("/scooters/closest", scooterHandler.GetClosestScooters)
		}
	}
}

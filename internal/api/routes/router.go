package routes

import (
	"net/http"
	"os"
	"path/filepath"

	"scootin-aboot/internal/api/handlers"
	"scootin-aboot/internal/api/middleware"
	"scootin-aboot/internal/services"
	"scootin-aboot/pkg/auth/apikey"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, apiKey string, tripService services.TripService, scooterService services.ScooterService) {
	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	scooterHandler := handlers.NewScooterHandler(tripService, scooterService)

	// Initialize API key validator
	apiKeyValidator := apikey.NewValidator(apiKey)

	// Documentation routes (no authentication required)
	router.GET("/docs", func(c *gin.Context) {
		// Serve Swagger UI
		swaggerUIPath := filepath.Join(".", "docs", "swagger-ui.html")
		if _, err := os.Stat(swaggerUIPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Documentation not found"})
			return
		}
		c.File(swaggerUIPath)
	})

	router.GET("/api-docs.yaml", func(c *gin.Context) {
		// Serve OpenAPI specification
		apiDocsPath := filepath.Join(".", "docs", "api", "openapi.yaml")
		if _, err := os.Stat(apiDocsPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "API documentation not found"})
			return
		}
		c.File(apiDocsPath)
	})

	// Serve OpenAPI component files
	router.Static("/paths", "./docs/api/paths")
	router.Static("/components", "./docs/api/components")
	router.Static("/info", "./docs/api/info")

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

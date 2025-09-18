package routes

import (
	"scootin-aboot/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine) {
	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Health check endpoint
		v1.GET("/health", healthHandler.HealthCheck)
	}
}

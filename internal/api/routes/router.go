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

func SetupRoutes(router *gin.Engine, apiKey string, tripService services.TripService, scooterService services.ScooterService) {
	healthHandler := handlers.NewHealthHandler()
	scooterHandler := handlers.NewScooterHandler(tripService, scooterService)

	apiKeyValidator := apikey.NewValidator(apiKey)

	router.GET("/docs", func(c *gin.Context) {
		swaggerUIPath := filepath.Join(".", "docs", "swagger-ui.html")
		if _, err := os.Stat(swaggerUIPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Documentation not found"})
			return
		}
		c.File(swaggerUIPath)
	})

	router.GET("/api-docs.yaml", func(c *gin.Context) {
		apiDocsPath := filepath.Join(".", "docs", "api", "openapi.yaml")
		if _, err := os.Stat(apiDocsPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "API documentation not found"})
			return
		}
		c.File(apiDocsPath)
	})

	router.Static("/paths", "./docs/api/paths")
	router.Static("/components", "./docs/api/components")
	router.Static("/info", "./docs/api/info")

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.HealthCheck)

		protected := v1.Group("")
		protected.Use(middleware.APIKeyMiddleware(apiKeyValidator))
		{
			protected.POST("/scooters/:id/trip/start", scooterHandler.StartTrip)
			protected.POST("/scooters/:id/trip/end", scooterHandler.EndTrip)
			protected.POST("/scooters/:id/location", scooterHandler.UpdateLocation)
			protected.GET("/scooters", scooterHandler.GetScooters)
			protected.GET("/scooters/:id", scooterHandler.GetScooter)
			protected.GET("/scooters/closest", scooterHandler.GetClosestScooters)
		}
	}
}

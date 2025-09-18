package main

import (
	"fmt"
	"log"

	"scootin-aboot/internal/api/routes"
	"scootin-aboot/internal/config"
	"scootin-aboot/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	if err := utils.InitLogger(cfg.LogLevel, cfg.LogFormat); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer utils.Sync()

	utils.Info("Starting Scootin' Aboot server",
		zap.String("host", cfg.ServerHost),
		zap.String("port", cfg.ServerPort),
		zap.String("log_level", cfg.LogLevel),
	)

	// Set Gin mode based on log level
	if cfg.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	address := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	utils.Info("Server starting", zap.String("address", address))

	if err := router.Run(address); err != nil {
		utils.Fatal("Failed to start server", zap.Error(err))
	}
}

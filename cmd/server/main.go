package main

import (
	"fmt"
	"log"

	"scootin-aboot/internal/api/middleware"
	"scootin-aboot/internal/api/routes"
	"scootin-aboot/internal/config"
	"scootin-aboot/internal/repository"
	"scootin-aboot/internal/services"
	"scootin-aboot/pkg/database"
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

	// Connect to database
	dsn := cfg.GetDatabaseDSN()
	gormDB, err := database.ConnectDatabase(dsn)
	if err != nil {
		utils.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Get underlying sql.DB for migrations
	sqlDB, err := gormDB.DB()
	if err != nil {
		utils.Fatal("Failed to get underlying sql.DB", zap.Error(err))
	}

	// Run database migrations
	migrationsPath, err := database.GetMigrationsPath()
	if err != nil {
		utils.Fatal("Failed to get migrations path", zap.Error(err))
	}

	if err := database.MigrateUp(sqlDB, migrationsPath); err != nil {
		utils.Fatal("Failed to run database migrations", zap.Error(err))
	}

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

	// Initialize repositories
	repo := repository.NewRepository(gormDB)

	// Initialize services
	tripService := services.NewTripService(
		repo.Trip(),
		repo.Scooter(),
		repo.User(),
		repo.LocationUpdate(),
	)

	scooterService := services.NewScooterService(
		repo.Scooter(),
		repo.Trip(),
	)

	// Create Gin router
	router := gin.New()

	// Add middleware (order matters!)
	router.Use(middleware.LoggingMiddleware())                // Custom logging with zap
	router.Use(middleware.RecoveryMiddleware())               // Custom recovery with structured logging
	router.Use(middleware.ErrorHandlerMiddleware())           // Error handling
	router.Use(middleware.ValidateJSON())                     // JSON validation
	router.Use(middleware.ValidateContentLength(1024 * 1024)) // 1MB max content length

	// Setup routes
	routes.SetupRoutes(router, cfg.APIKey, tripService, scooterService)

	// Start server
	address := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	utils.Info("Server starting", zap.String("address", address))

	if err := router.Run(address); err != nil {
		utils.Fatal("Failed to start server", zap.Error(err))
	}
}

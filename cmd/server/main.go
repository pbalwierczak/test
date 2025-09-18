package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Create a separate connection for migrations to avoid closing the main connection
	migrationDB, err := database.ConnectDatabase(dsn)
	if err != nil {
		utils.Fatal("Failed to connect to database for migrations", zap.Error(err))
	}

	// Get underlying sql.DB for migrations
	sqlDB, err := migrationDB.DB()
	if err != nil {
		utils.Fatal("Failed to get underlying sql.DB for migrations", zap.Error(err))
	}

	// Run database migrations
	migrationsPath, err := database.GetMigrationsPath()
	if err != nil {
		utils.Fatal("Failed to get migrations path", zap.Error(err))
	}

	if err := database.MigrateUp(sqlDB, migrationsPath); err != nil {
		utils.Fatal("Failed to run database migrations", zap.Error(err))
	}

	// Close the migration connection after migrations are done
	if err := sqlDB.Close(); err != nil {
		utils.Error("Failed to close migration database connection", zap.Error(err))
	}

	// Start database health check (every 30 seconds)
	stopHealthCheck := database.StartHealthCheck(gormDB, 30*time.Second)

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

	// Start server with graceful shutdown
	address := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	utils.Info("Server starting", zap.String("address", address))

	// Create a server instance for graceful shutdown
	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Info("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utils.Fatal("Server forced to shutdown", zap.Error(err))
	}

	// Stop database health check
	stopHealthCheck()

	// Close database connection gracefully
	mainSqlDB, err := gormDB.DB()
	if err != nil {
		utils.Error("Failed to get main database connection for closing", zap.Error(err))
	} else {
		if err := mainSqlDB.Close(); err != nil {
			utils.Error("Failed to close database connection", zap.Error(err))
		}
	}

	utils.Info("Server exited")
}

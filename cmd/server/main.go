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
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := utils.InitLogger(cfg.LogLevel, cfg.LogFormat); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer utils.Sync()

	dsn := cfg.GetDatabaseDSN()
	gormDB, err := database.ConnectDatabase(dsn)
	if err != nil {
		utils.Fatal("Failed to connect to database", zap.Error(err))
	}

	migrationDB, err := database.ConnectDatabase(dsn)
	if err != nil {
		utils.Fatal("Failed to connect to database for migrations", zap.Error(err))
	}

	sqlDB, err := migrationDB.DB()
	if err != nil {
		utils.Fatal("Failed to get underlying sql.DB for migrations", zap.Error(err))
	}

	migrationsPath, err := database.GetMigrationsPath()
	if err != nil {
		utils.Fatal("Failed to get migrations path", zap.Error(err))
	}

	if err := database.MigrateUp(sqlDB, migrationsPath); err != nil {
		utils.Fatal("Failed to run database migrations", zap.Error(err))
	}

	if err := sqlDB.Close(); err != nil {
		utils.Error("Failed to close migration database connection", zap.Error(err))
	}

	stopHealthCheck := database.StartHealthCheck(gormDB, 30*time.Second)

	utils.Info("Starting Scootin' Aboot server",
		zap.String("host", cfg.ServerHost),
		zap.String("port", cfg.ServerPort),
		zap.String("log_level", cfg.LogLevel),
	)

	if cfg.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	repo := repository.NewRepository(gormDB)

	tripService := services.NewTripService(
		repo.Trip(),
		repo.Scooter(),
		repo.User(),
		repo.LocationUpdate(),
	)

	scooterService := services.NewScooterService(
		repo.Scooter(),
		repo.Trip(),
		repo.LocationUpdate(),
	)

	router := gin.New()

	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.ValidateJSON())
	router.Use(middleware.ValidateContentLength(1024 * 1024))

	routes.SetupRoutes(router, cfg.APIKey, tripService, scooterService)

	address := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	utils.Info("Server starting", zap.String("address", address))

	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utils.Fatal("Server forced to shutdown", zap.Error(err))
	}

	stopHealthCheck()

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

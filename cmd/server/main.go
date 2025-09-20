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
	"scootin-aboot/internal/kafka"
	"scootin-aboot/internal/repository"
	"scootin-aboot/internal/services"
	"scootin-aboot/pkg/database"
	"scootin-aboot/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := logger.InitLogger(cfg.LogLevel, cfg.LogFormat); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	dsn := cfg.GetDatabaseDSN()
	dbURL := cfg.GetDatabaseURL()

	migrationsPath, err := database.GetMigrationsPath()
	if err != nil {
		logger.Fatal("Failed to get migrations path", logger.ErrorField(err))
	}

	if err := database.MigrateUp(dbURL, migrationsPath); err != nil {
		logger.Fatal("Failed to run database migrations", logger.ErrorField(err))
	}

	sqlDB, err := database.ConnectDatabase(dsn)
	if err != nil {
		logger.Fatal("Failed to connect to database", logger.ErrorField(err))
	}

	stopHealthCheck := database.StartHealthCheck(sqlDB, 30*time.Second)

	logger.Info("Starting Scootin' Aboot server",
		logger.String("host", cfg.ServerHost),
		logger.String("port", cfg.ServerPort),
		logger.String("log_level", cfg.LogLevel),
	)

	if cfg.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	repo := repository.NewRepository(sqlDB)

	tripService := services.NewTripService(
		repo.Trip(),
		repo.Scooter(),
		repo.User(),
		repo.LocationUpdate(),
		repo.UnitOfWork(),
	)

	scooterService := services.NewScooterService(
		repo.Scooter(),
		repo.Trip(),
		repo.LocationUpdate(),
		repo.UnitOfWork(),
	)

	router := gin.New()

	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.ValidateJSON())
	router.Use(middleware.ValidateContentLength(1024 * 1024))

	routes.SetupRoutes(router, cfg.APIKey, scooterService)

	kafkaConsumer, err := kafka.NewEventConsumer(&cfg.KafkaConfig, tripService, scooterService)
	if err != nil {
		logger.Fatal("Failed to create Kafka consumer", logger.ErrorField(err))
	}

	if err := kafkaConsumer.Start(); err != nil {
		logger.Fatal("Failed to start Kafka consumer", logger.ErrorField(err))
	}
	logger.Info("Kafka consumer started")

	address := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	logger.Info("Server starting", logger.String("address", address))

	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", logger.ErrorField(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", logger.ErrorField(err))
	}

	stopHealthCheck()

	kafkaConsumer.Stop()
	logger.Info("Kafka consumer stopped")

	if err := sqlDB.Close(); err != nil {
		logger.Error("Failed to close database connection", logger.ErrorField(err))
	}

	logger.Info("Server exited")
}

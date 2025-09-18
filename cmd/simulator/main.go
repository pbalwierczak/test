package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"scootin-aboot/internal/config"
	"scootin-aboot/pkg/utils"

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

	utils.Info("Starting Scootin' Aboot simulator",
		zap.Int("scooters", cfg.SimulatorScooters),
		zap.Int("users", cfg.SimulatorUsers),
		zap.String("server_url", cfg.SimulatorServerURL),
	)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan
	utils.Info("Received shutdown signal, stopping simulator...")
}

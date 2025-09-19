package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"scootin-aboot/internal/config"
	"scootin-aboot/pkg/simulator"
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

	// Log comprehensive configuration
	utils.Info("Simulator Configuration",
		zap.Int("scooters", cfg.SimulatorScooters),
		zap.Int("users", cfg.SimulatorUsers),
		zap.String("server_url", cfg.SimulatorServerURL),
		zap.Int("speed", cfg.SimulatorSpeed),
		zap.Int("trip_duration_min", cfg.SimulatorTripDurationMin),
		zap.Int("trip_duration_max", cfg.SimulatorTripDurationMax),
		zap.Int("rest_min", cfg.SimulatorRestMin),
		zap.Int("rest_max", cfg.SimulatorRestMax),
		zap.String("log_level", cfg.LogLevel),
		zap.String("log_format", cfg.LogFormat),
	)

	utils.Info("Starting Scootin' Aboot simulator")

	// Create simulator
	sim := simulator.NewSimulator(cfg)

	// Start simulation
	if err := sim.Start(); err != nil {
		utils.Fatal("Failed to start simulator", zap.Error(err))
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan
	utils.Info("Received shutdown signal, stopping simulator...")

	// Stop simulation
	sim.Stop()
}

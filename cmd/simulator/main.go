package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"scootin-aboot/internal/config"
	"scootin-aboot/pkg/logger"
	"scootin-aboot/pkg/simulator"
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

	logger.Info("Simulator Configuration",
		logger.Int("scooters", cfg.SimulatorScooters),
		logger.Int("users", cfg.SimulatorUsers),
		logger.String("server_url", cfg.SimulatorServerURL),
		logger.Int("speed", cfg.SimulatorSpeed),
		logger.Int("trip_duration_min", cfg.SimulatorTripDurationMin),
		logger.Int("trip_duration_max", cfg.SimulatorTripDurationMax),
		logger.Int("rest_min", cfg.SimulatorRestMin),
		logger.Int("rest_max", cfg.SimulatorRestMax),
		logger.String("log_level", cfg.LogLevel),
		logger.String("log_format", cfg.LogFormat),
	)

	logger.Info("Starting Scootin' Aboot simulator")

	sim := simulator.NewSimulator(cfg)

	if err := sim.Start(); err != nil {
		logger.Fatal("Failed to start simulator", logger.ErrorField(err))
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("Received shutdown signal, stopping simulator...")

	sim.Stop()
}

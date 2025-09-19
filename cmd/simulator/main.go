package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"scootin-aboot/internal/config"
	"scootin-aboot/pkg/simulator"
	"scootin-aboot/pkg/utils"
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

	utils.Info("Simulator Configuration",
		utils.Int("scooters", cfg.SimulatorScooters),
		utils.Int("users", cfg.SimulatorUsers),
		utils.String("server_url", cfg.SimulatorServerURL),
		utils.Int("speed", cfg.SimulatorSpeed),
		utils.Int("trip_duration_min", cfg.SimulatorTripDurationMin),
		utils.Int("trip_duration_max", cfg.SimulatorTripDurationMax),
		utils.Int("rest_min", cfg.SimulatorRestMin),
		utils.Int("rest_max", cfg.SimulatorRestMax),
		utils.String("log_level", cfg.LogLevel),
		utils.String("log_format", cfg.LogFormat),
	)

	utils.Info("Starting Scootin' Aboot simulator")

	sim := simulator.NewSimulator(cfg)

	if err := sim.Start(); err != nil {
		utils.Fatal("Failed to start simulator", utils.ErrorField(err))
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	utils.Info("Received shutdown signal, stopping simulator...")

	sim.Stop()
}

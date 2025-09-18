# Scootin' Aboot - Electric Scooter Management System
# Makefile for development and build automation

.PHONY: help build server simulator test clean deps

# Default target
help:
	@echo "Scootin' Aboot - Available targets:"
	@echo "  build      - Build both server and simulator"
	@echo "  server     - Run the main API server"
	@echo "  simulator  - Run the simulation program"
	@echo "  test       - Run all tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  deps       - Download dependencies"
	@echo "  help       - Show this help message"

# Build targets
build: deps
	@echo "Building server..."
	@go build -o bin/server ./cmd/server
	@echo "Building simulator..."
	@go build -o bin/simulator ./cmd/simulator
	@echo "Build complete!"

# Run targets
server: deps
	@echo "Starting server..."
	@go run ./cmd/server

simulator: deps
	@echo "Starting simulator..."
	@go run ./cmd/simulator

# Test targets
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean targets
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@go clean

# Dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

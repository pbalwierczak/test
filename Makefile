# Scootin' Aboot - Electric Scooter Management System
# Makefile for development and build automation

.PHONY: app simulator db build clean test help _build _server _simulator _test _clean _deps docs docs-validate docs-clean

# Default target
app:
	@echo "Building and starting app with database..."
	@docker-compose up --build

# Help target
help:
	@echo "Scootin' Aboot - Available targets:"
	@echo "  app          - Build and run app with database (default)"
	@echo "  simulator    - Build and run simulator"
	@echo "  db           - Run database only"
	@echo "  build        - Build all Docker images"
	@echo "  clean        - Clean up Docker containers and images"
	@echo "  test         - Run tests in Docker container"
	@echo ""
	@echo "Local development targets (prefixed with _):"
	@echo "  _build      - Build both server and simulator locally"
	@echo "  _server     - Run the main API server locally"
	@echo "  _simulator  - Run the simulation program locally"
	@echo "  _test       - Run all tests locally"
	@echo "  _clean      - Clean build artifacts"
	@echo "  _deps       - Download dependencies"
	@echo "  seed        - Load seed data into database"
	@echo "  seed-reset  - Reset database and reload seeds"
	@echo ""
	@echo "Documentation targets:"
	@echo "  docs          - Show documentation help"
	@echo "  docs-validate - Validate OpenAPI specification using Docker"
	@echo "  docs-clean    - Clean documentation artifacts"
	@echo ""
	@echo "API Documentation:"
	@echo "  Swagger UI is available at: http://localhost:8080/docs"
	@echo "  OpenAPI spec is available at: http://localhost:8080/api-docs.yaml"
	@echo "  help        - Show this help message"

# Local development targets (prefixed with _)
_build: _deps
	@echo "Building server..."
	@go build -o bin/server ./cmd/server
	@echo "Building simulator..."
	@go build -o bin/simulator ./cmd/simulator
	@echo "Build complete!"

_server: _deps
	@echo "Starting server..."
	@go run ./cmd/server

_simulator: _deps
	@echo "Starting simulator..."
	@go run ./cmd/simulator

_test:
	@echo "Running tests..."
	@go test -v ./...

_clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@go clean

_deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Docker targets (primary workflow)
simulator:
	@echo "Building and starting simulator..."
	@docker-compose -f docker-compose.simulator.yml up --build

db:
	@echo "Starting database only..."
	@docker-compose up postgres -d

build:
	@echo "Building all Docker images..."
	@docker build -t scootin-app .
	@docker build -f Dockerfile.simulator -t scootin-simulator .

clean:
	@echo "Cleaning up Docker containers and images..."
	@docker-compose down -v
	@docker-compose -f docker-compose.simulator.yml down
	@docker rmi scootin-app scootin-simulator 2>/dev/null || true
	@docker system prune -f

test:
	@echo "Running tests in Docker container..."
	@docker build --target test -t scootin-test .
	@docker run --rm scootin-test go test -v ./...

# Database seed commands
seed:
	@echo "Loading seed data into database..."
	@docker exec -i scootin-postgres psql -U postgres -d scootin_aboot < seeds/users.sql
	@docker exec -i scootin-postgres psql -U postgres -d scootin_aboot < seeds/scooters.sql
	@docker exec -i scootin-postgres psql -U postgres -d scootin_aboot < seeds/sample_trips.sql
	@echo "Seed data loaded successfully!"

seed-reset:
	@echo "Resetting database and loading seeds..."
	@docker-compose down -v
	@docker-compose up postgres -d
	@echo "Waiting for database to be ready..."
	@sleep 10
	@$(MAKE) seed
	@echo "Database reset and seeded successfully!"

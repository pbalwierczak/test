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
	@echo "  simulator    - Build and run simulator (Docker)"
	@echo "  simulator-test - Test simulator connectivity to external app"
	@echo "  simulator-run - Run simulator locally with logging"
	@echo "  simulator-stop - Stop running simulator"
	@echo "  simulator-tail-logs - Follow simulator logs"
	@echo "  simulator-clean-logs - Clean log files"
	@echo "  dev          - Run both server and simulator locally"
	@echo "  dev-stop     - Stop development environment"
	@echo "  db           - Run database only"
	@echo "  build        - Build all Docker images"
	@echo "  clean        - Clean up Docker containers and images"
	@echo "  test         - Run tests in Docker container"
	@echo ""
	@echo "Local development targets (prefixed with _):"
	@echo "  _build      - Build both server and simulator locally"
	@echo "  _server     - Run the main API server locally"
	@echo "  _simulator  - Run the simulation program locally"
	@echo "  _simulator-log - Run simulator with logging to file"
	@echo "  _stop-simulator - Stop running simulator"
	@echo "  _tail-logs   - Follow simulator logs in real-time"
	@echo "  _dev         - Run both server and simulator for development"
	@echo "  _stop-dev    - Stop development environment"
	@echo "  _test       - Run all tests locally"
	@echo "  _clean      - Clean build artifacts"
	@echo "  _clean-logs - Clean log files"
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

simulator-clean-logs:
	@echo "Cleaning log files..."
	@rm -f simulator.log
	@rm -f server.log
	@rm -f *.log
	@echo "Log files cleaned!"

simulator-tail-logs:
	@echo "Following simulator logs (Ctrl+C to stop)..."
	@tail -f simulator.log

simulator-run:
	@echo "Starting simulator with logging..."
	@go run ./cmd/simulator > simulator.log 2>&1 &
	@echo "Simulator started in background. Logs are being written to simulator.log"
	@echo "Use 'make simulator-stop' to stop the simulator"
	@echo "Use 'make simulator-tail-logs' to follow the logs"

simulator-stop:
	@echo "Stopping simulator..."
	@pkill -f "go run ./cmd/simulator" || pkill -f "./bin/simulator" || echo "No simulator process found"
	@echo "Simulator stopped"

dev:
	@echo "Starting development environment..."
	@echo "Starting server in background..."
	@go run ./cmd/server > server.log 2>&1 &
	@sleep 3
	@echo "Starting simulator in background..."
	@go run ./cmd/simulator > simulator.log 2>&1 &
	@echo "Development environment started!"
	@echo "Server logs: server.log"
	@echo "Simulator logs: simulator.log"
	@echo "Use 'make dev-stop' to stop both services"
	@echo "Use 'make simulator-tail-logs' to follow simulator logs"

dev-stop:
	@echo "Stopping development environment..."
	@pkill -f "go run ./cmd/server" || pkill -f "./bin/server" || echo "No server process found"
	@pkill -f "go run ./cmd/simulator" || pkill -f "./bin/simulator" || echo "No simulator process found"
	@echo "Development environment stopped"

_deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Docker targets (primary workflow)
simulator:
	@echo "Building and starting simulator..."
	@echo "Note: Make sure the main app is running first with 'make app'"
	@docker-compose -f docker-compose.simulator.yml up --build

simulator-test:
	@echo "Testing simulator connectivity..."
	@if ! docker ps | grep -q "scootin-app"; then \
		echo "❌ Error: scootin-app container is not running!"; \
		echo "Please start the main application first with: make app"; \
		exit 1; \
	fi
	@echo "✅ scootin-app container is running"
	@if curl -s -f "http://localhost:8080/health" > /dev/null; then \
		echo "✅ API is accessible from host"; \
	else \
		echo "❌ API is not accessible from host"; \
		echo "Make sure the main application is running and accessible on port 8080"; \
		exit 1; \
	fi
	@echo "🎉 Simulator connectivity test passed!"

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

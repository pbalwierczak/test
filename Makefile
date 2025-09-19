.PHONY: app simulator db build clean test help _build _server _simulator _test _clean _deps docs docs-validate docs-clean

app:
	@echo "Building and starting app with database..."
	@echo "Cleaning up any existing containers..."
	@docker-compose down -v 2>/dev/null || true
	@echo "Starting services..."
	@docker-compose up --build -d
	@echo "Waiting for services to be ready..."
	@echo "Checking service status..."
	@for i in $$(seq 1 60); do \
		if docker-compose ps | grep -q "scootin-app.*Up"; then \
			break; \
		fi; \
		echo "Waiting for services... ($$i/60)"; \
		sleep 2; \
	done
	@if ! docker-compose ps | grep -q "scootin-app.*Up"; then \
		echo "❌ Services failed to start properly"; \
		docker-compose logs --tail=50; \
		exit 1; \
	fi
	@echo "✅ All services are running!"
	@echo "📊 Service status:"
	@docker-compose ps
	@echo ""
	@echo "🌐 Application available at: http://localhost:8080"
	@echo "📚 API docs available at: http://localhost:8080/docs"
	@echo "💚 Health check: http://localhost:8080/api/v1/health"
	@echo ""
	@echo "To view logs: docker-compose logs -f"
	@echo "To stop: docker-compose down"

help:
	@echo "Scootin' Aboot - Available targets:"
	@echo "  app          - Build and run app with database (default)"
	@echo "  simulator    - Build and run simulator (Docker)"
	@echo "  simulator-tail-logs - Follow simulator logs"
	@echo "  simulator-clean-logs - Clean log files"
	@echo "  build        - Build all Docker images"
	@echo "  clean        - Clean up Docker containers and images"
	@echo "  test         - Run tests in Docker container"
	@echo "  logs         - Follow application logs"
	@echo "  status       - Show service status and health"
	@echo "  stop         - Stop all services"
	@echo "  restart      - Restart all services"
	@echo ""
	@echo "Local development targets (prefixed with _):"
	@echo "  _build      - Build both server and simulator locally"
	@echo "  _test       - Run all tests locally"
	@echo "  _clean      - Clean build artifacts"
	@echo "  _deps       - Download dependencies"
	@echo "  seed        - Truncate and reload seed data into database"
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

_build: _deps
	@echo "Building server..."
	@go build -o bin/server ./cmd/server
	@echo "Building simulator..."
	@go build -o bin/simulator ./cmd/simulator
	@echo "Build complete!"

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

_deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

simulator:
	@echo "Building and starting simulator..."
	@if ! docker network ls | grep -q "scootin-aboot-app_scootin-network"; then \
		echo "❌ Error: scootin-network does not exist!"; \
		echo "Please start the main application first with: make app"; \
		exit 1; \
	fi
	@if ! docker ps | grep -q "scootin-app"; then \
		echo "❌ Error: scootin-app container is not running!"; \
		echo "Please start the main application first with: make app"; \
		exit 1; \
	fi
	@echo "✅ Main app is running, starting simulator..."
	@docker-compose -f docker-compose.simulator.yml up --build

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

logs:
	@echo "Following application logs (Ctrl+C to stop)..."
	@docker-compose logs -f

status:
	@echo "Service status:"
	@docker-compose ps
	@echo ""
	@echo "Health checks:"
	@docker-compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}"

stop:
	@echo "Stopping all services..."
	@docker-compose down

restart:
	@echo "Restarting services..."
	@docker-compose restart

test:
	@echo "Running tests in Docker container..."
	@docker build --target test -t scootin-test .
	@docker run --rm scootin-test go test -v ./...

seed:
	@echo "⚠️  WARNING: This will TRUNCATE all tables and reload seed data!"
	@echo "Tables that will be cleared: users, scooters, trips, location_updates"
	@echo ""
	@read -p "Are you sure you want to continue? (y/N): " confirm && [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ] || (echo "Operation cancelled." && exit 1)
	@echo "Loading seed data into database..."
	@docker exec -i scootin-postgres psql -U postgres -d scootin_aboot < seeds/users.sql
	@docker exec -i scootin-postgres psql -U postgres -d scootin_aboot < seeds/scooters.sql
	@docker exec -i scootin-postgres psql -U postgres -d scootin_aboot < seeds/sample_trips.sql
	@echo "Seed data loaded successfully!"


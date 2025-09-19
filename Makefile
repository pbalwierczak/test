.PHONY: app simulator start-app start-sim start-simulator logs-app logs-sim logs-simulator kill-app kill-sim kill-simulator kill-all status clean test seed help

# Default target
.DEFAULT_GOAL := help

help:
	@echo "🚀 Scootin' Aboot - Quick Start Commands:"
	@echo ""
	@echo "  make start-app     - Start app with database in background"
	@echo "  make start-sim     - Start simulator in background (alias: start-simulator)"
	@echo "  make logs-app      - View app logs (Ctrl+C to exit)"
	@echo "  make logs-sim      - View simulator logs (Ctrl+C to exit, alias: logs-simulator)"
	@echo "  make kill-app      - Stop app and database"
	@echo "  make kill-sim      - Stop simulator (alias: kill-simulator)"
	@echo "  make kill-all      - Stop everything"
	@echo "  make status        - Show running services"
	@echo ""
	@echo "🌐 Once running:"
	@echo "  App: http://localhost:8080"
	@echo "  API docs: http://localhost:8080/docs"
	@echo "  Health: http://localhost:8080/api/v1/health"
	@echo ""
	@echo "📋 Other useful commands:"
	@echo "  make seed          - Load sample data into database"
	@echo "  make clean         - Clean up everything"
	@echo "  make test          - Run tests"

start-app:
	@echo "🚀 Starting app with database in background..."
	@echo "Cleaning up any existing containers..."
	@docker-compose down -v 2>/dev/null || true
	@echo "Starting services..."
	@docker-compose up --build -d
	@echo "Waiting for services to be ready..."
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
	@echo "✅ App is running in background!"
	@echo "🌐 Available at: http://localhost:8080"
	@echo "📚 API docs: http://localhost:8080/docs"
	@echo ""
	@echo "💡 Next steps:"
	@echo "  make logs-app      - View app logs"
	@echo "  make start-sim     - Start simulator"
	@echo "  make seed          - Load sample data"
	@echo "  make status        - Check service status"
	@echo "  make kill-app      - Stop app"

start-sim: start-simulator
start-simulator:
	@echo "🎮 Starting simulator in background..."
	@if ! docker network ls | grep -q "scootin-aboot-app_scootin-network"; then \
		echo "❌ Error: Network not found! Please start the app first: make start-app"; \
		exit 1; \
	fi
	@if ! docker ps | grep -q "scootin-app"; then \
		echo "❌ Error: App not running! Please start the app first: make start-app"; \
		exit 1; \
	fi
	@echo "✅ App is running, starting simulator..."
	@docker-compose -f docker-compose.simulator.yml up --build -d
	@echo "✅ Simulator is running in background!"
	@echo ""
	@echo "💡 Next steps:"
	@echo "  make logs-sim      - View simulator logs"
	@echo "  make status        - Check all services"
	@echo "  make kill-sim      - Stop simulator"

logs-app:
	@echo "📋 Following app logs (Ctrl+C to exit)..."
	@docker-compose logs -f scootin-app

logs-sim: logs-simulator
logs-simulator:
	@echo "📋 Following simulator logs (Ctrl+C to exit)..."
	@docker-compose -f docker-compose.simulator.yml logs -f scootin-simulator

kill-app:
	@echo "🛑 Stopping app and database..."
	@docker-compose down
	@echo "✅ App stopped!"

kill-sim: kill-simulator
kill-simulator:
	@echo "🛑 Stopping simulator..."
	@docker-compose -f docker-compose.simulator.yml down
	@echo "✅ Simulator stopped!"

kill-all:
	@echo "🛑 Stopping everything..."
	@docker-compose down
	@docker-compose -f docker-compose.simulator.yml down
	@echo "✅ All services stopped!"

status:
	@echo "📊 Service Status:"
	@echo ""
	@echo "Main App Services:"
	@docker-compose ps --services 2>/dev/null | while read service; do \
		if [ "$$service" != "scootin-simulator" ]; then \
			docker-compose ps $$service 2>/dev/null | tail -n +2; \
		fi; \
	done || echo "  No main app services running"
	@echo ""
	@echo "Simulator Services:"
	@docker-compose -f docker-compose.simulator.yml ps --services 2>/dev/null | while read service; do \
		docker-compose -f docker-compose.simulator.yml ps $$service 2>/dev/null | tail -n +2; \
	done || echo "  No simulator services running"
	@echo ""
	@if docker ps | grep -q "scootin-app"; then \
		echo "🌐 App is running at: http://localhost:8080"; \
		echo "📚 API docs: http://localhost:8080/docs"; \
	fi

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

clean:
	@echo "🧹 Cleaning up everything..."
	@docker-compose down -v
	@docker-compose -f docker-compose.simulator.yml down
	@docker rmi scootin-app scootin-simulator 2>/dev/null || true
	@docker system prune -f
	@echo "✅ Cleanup complete!"

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

# Convenience aliases
app: start-app
simulator: start-simulator
logs: logs-app
stop: kill-all


# Scootin' Aboot - Electric Scooter Management System

A backend service for managing electric scooters in Ottawa and Montreal, providing REST API for scooter event collection and mobile client reporting.

## Features

- **Scooter Management**: Track scooter status, location, and trip lifecycle
- **Real-time Location Updates**: Periodic GPS updates during trips
- **Geographic Search**: Find scooters by location and status
- **Simulation**: Built-in simulator for testing with fake clients and scooters
- **REST API**: Clean, documented API endpoints
- **Docker Support**: Containerized deployment

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL (for production)
- Docker & Docker Compose (optional)

### Development Setup

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd scootin-aboot-app
   cp .env.example .env
   # Edit .env with your configuration
   ```

2. **Install dependencies**:
   ```bash
   make deps
   ```

3. **Run the server**:
   ```bash
   make server
   ```

4. **Run the simulator** (in another terminal):
   ```bash
   make simulator
   ```

### Available Commands

- `make build` - Build both server and simulator
- `make server` - Run the main API server
- `make simulator` - Run the simulation program
- `make test` - Run all tests
- `make clean` - Clean build artifacts
- `make deps` - Download dependencies

## API Endpoints

### Health Check
- `GET /api/v1/health` - Service health status

### Scooter Management
- `POST /api/v1/scooters/{id}/trip/start` - Start a trip
- `POST /api/v1/scooters/{id}/trip/end` - End a trip
- `POST /api/v1/scooters/{id}/location` - Update location
- `GET /api/v1/scooters` - Query scooters with filters
- `GET /api/v1/scooters/{id}` - Get specific scooter details
- `GET /api/v1/scooters/closest` - Find closest scooters

## Configuration

Configuration is managed through environment variables. See `.env.example` for all available options.

Key configuration areas:
- **API**: Server host/port, API key
- **Database**: PostgreSQL connection settings
- **Simulator**: Number of scooters/users, behavior parameters
- **Geographic**: City coordinates and search radius
- **Logging**: Log level and format

## Project Structure

```
scootin-aboot/
├── cmd/                    # Application entry points
│   ├── server/            # Main API server
│   └── simulator/         # Simulation program
├── internal/              # Private application code
│   ├── api/              # HTTP handlers and middleware
│   ├── config/           # Configuration management
│   ├── models/           # Data models
│   ├── repository/       # Data access layer
│   └── services/         # Business logic
├── pkg/                   # Public packages
│   ├── auth/             # Authentication
│   ├── database/         # Database utilities
│   ├── simulator/        # Simulation logic
│   └── utils/            # Common utilities
├── migrations/            # Database migrations
├── seeds/                 # Seed data
└── docs/                  # Documentation
```

## Development

This project follows Go best practices:
- Clean architecture with separated concerns
- Comprehensive testing
- Structured logging
- Environment-based configuration
- Docker containerization

## License

[Add your license here]

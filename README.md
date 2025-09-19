# Scootin' Aboot - Electric Scooter Management System

A comprehensive backend service for managing electric scooters in Ottawa and Montreal, featuring both REST API and event-driven architecture for scooter event collection and mobile client reporting.

## Features

- **Scooter Management**: Track scooter status, location, and complete trip lifecycle
- **Real-time Location Updates**: Periodic GPS updates during trips with Kafka event streaming
- **Geographic Search**: Advanced location-based filtering and closest scooter discovery
- **Event-Driven Architecture**: Kafka-based communication for real-time processing
- **Simulation System**: Built-in simulator with realistic scooter and user behavior
- **Comprehensive API**: Fully documented REST API with OpenAPI 3.0 specification
- **Containerized Deployment**: Complete Docker Compose setup with PostgreSQL and Kafka
- **Authentication**: API key-based security for all endpoints
- **Database Management**: Automated migrations and seed data

## Quick Start

### Prerequisites

- Go 1.25+
- Docker & Docker Compose (recommended)
- PostgreSQL 15+ (included in Docker setup)
- Apache Kafka 7.4+ (included in Docker setup)

### Development Setup

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd scootin-aboot-app
   cp .env.example .env
   # Edit .env with your configuration
   ```

2. **Start the application** (includes database and Kafka):
   ```bash
   make start-app
   ```

3. **Start the simulator** (in another terminal):
   ```bash
   make start-sim
   ```

4. **View logs**:
   ```bash
   make logs-app    # View application logs
   make logs-sim    # View simulator logs
   ```

### Available Commands

**Quick Start:**
- `make start-app` - Start app with database and Kafka in background
- `make start-sim` - Start simulator in background
- `make logs-app` - View app logs (Ctrl+C to exit)
- `make logs-sim` - View simulator logs (Ctrl+C to exit)
- `make status` - Show running services
- `make kill-all` - Stop everything

**Development:**
- `make test` - Run all tests in Docker container
- `make seed` - Load sample data into database
- `make clean` - Clean up everything (containers, images, volumes)

**Direct Access:**
- App: http://localhost:8080
- API docs: http://localhost:8080/docs
- Health check: http://localhost:8080/api/v1/health

## API Endpoints

All endpoints require API key authentication via the `Authorization` header, except for the health check endpoint.

### System
- `GET /api/v1/health` - Service health status (public endpoint)

### Scooter Management
- `GET /api/v1/scooters` - List scooters with geographic and status filtering
  - Query parameters: `status`, `min_lat`, `max_lat`, `min_lng`, `max_lng`, `limit`, `offset`
- `GET /api/v1/scooters/{id}` - Get specific scooter details
- `GET /api/v1/scooters/closest` - Find closest scooters by location
  - Query parameters: `lat`, `lng`, `radius`, `status`, `limit`


### API Documentation
- Interactive API docs: http://localhost:8080/docs
- OpenAPI specification: http://localhost:8080/api/v1/openapi.json

## Event-Driven Architecture

The system uses Kafka event-driven communication for real-time processing and scalability. The simulator communicates with the server exclusively through Kafka events.

### Kafka Events

The simulator publishes events to Kafka topics that are consumed by the server:

- **Trip Started**: `scooter.trip.started` - When a user begins a trip
- **Trip Ended**: `scooter.trip.ended` - When a trip is completed
- **Location Updated**: `scooter.location.updated` - Periodic location updates during trips

### Event Flow

1. **Simulator** → Publishes events to Kafka topics
2. **Kafka** → Stores and distributes events
3. **Server** → Consumes events and processes them through existing services

### Benefits of Event-Driven Architecture

- **Decoupling**: Simulator and server operate independently
- **Scalability**: Multiple consumers can process events
- **Reliability**: Events are persisted and can be replayed
- **Observability**: Complete event flow visibility
- **Flexibility**: Easy to add new event consumers

## Configuration

Configuration is managed through environment variables. Create a `.env` file based on the example configuration.

### Key Configuration Areas

**API Server:**
- `SERVER_PORT`: HTTP server port (default: 8080)
- `API_KEY`: Static API key for authentication
- `LOG_LEVEL`: Logging level (debug, info, warn, error)

**Database:**
- `DB_HOST`: PostgreSQL host
- `DB_PORT`: PostgreSQL port (default: 5432)
- `DB_NAME`: Database name
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_SSLMODE`: SSL mode (disable for development)

**Kafka:**
- `KAFKA_BROKERS`: Kafka broker addresses
- `KAFKA_CLIENT_ID`: Client identifier
- `KAFKA_SECURITY_PROTOCOL`: Security protocol (PLAINTEXT for development)

**Simulator:**
- `SIMULATOR_SCOOTERS`: Number of scooters to simulate
- `SIMULATOR_USERS`: Number of users to simulate
- `SIMULATOR_INTERVAL`: Update interval in seconds

**Geographic:**
- `CITY_CENTER_LAT`: City center latitude
- `CITY_CENTER_LNG`: City center longitude
- `SEARCH_RADIUS`: Default search radius in meters

## Project Structure

```
scootin-aboot-app/
├── cmd/                    # Application entry points
│   ├── server/            # Main API server
│   └── simulator/         # Simulation program
├── internal/              # Private application code
│   ├── api/              # HTTP handlers, middleware, and routes
│   │   ├── handlers/     # Request handlers
│   │   ├── middleware/   # Auth, validation, logging
│   │   └── routes/       # Route definitions
│   ├── config/           # Configuration management
│   ├── kafka/            # Kafka consumer implementation
│   ├── models/           # Domain models and business logic
│   ├── repository/       # Data access layer (GORM)
│   └── services/         # Business logic services
├── pkg/                   # Public packages
│   ├── auth/             # API key authentication
│   ├── database/         # Database connection and migrations
│   ├── kafka/            # Kafka producer and events
│   ├── logger/           # Structured logging
│   ├── simulator/        # Simulation logic and movement
│   └── validation/       # Input validation utilities
├── migrations/            # Database schema migrations
├── seeds/                 # Sample data for development
├── docs/                  # API documentation (OpenAPI)
├── bin/                   # Built binaries
└── docker-compose*.yml    # Container orchestration
```

## Docker Services

The application runs in a containerized environment with the following services:

### Main Services (`docker-compose.yml`)
- **scootin-app**: Main API server (Go application)
- **postgres**: PostgreSQL 15 database with health checks
- **kafka**: Apache Kafka 7.4 with Zookeeper
- **zookeeper**: Kafka coordination service

### Simulator Services (`docker-compose.simulator.yml`)
- **scootin-simulator**: Simulation program for testing

### Key Features
- **Health Checks**: All services include proper health monitoring
- **Volume Persistence**: Database and Kafka data persist between restarts
- **Network Isolation**: Services communicate through dedicated Docker network
- **Environment Configuration**: All settings configurable via environment variables

## Technology Stack

### Backend
- **Go 1.25+**: Modern Go with generics and performance optimizations
- **Gin**: High-performance HTTP web framework
- **GORM**: Feature-rich ORM for database operations
- **PostgreSQL 15**: ACID-compliant relational database
- **Apache Kafka 7.4**: Event streaming platform for real-time data

### Infrastructure
- **Docker & Docker Compose**: Containerization and orchestration
- **Zookeeper**: Kafka coordination service
- **Health Checks**: Built-in service monitoring

### Development & Testing
- **Testify**: Comprehensive testing framework
- **Zap**: High-performance structured logging
- **OpenAPI 3.0**: API documentation and specification
- **Golang Migrate**: Database schema management

### Authentication & Security
- **API Key Authentication**: Simple, secure authentication
- **Input Validation**: Comprehensive request validation
- **CORS Support**: Cross-origin resource sharing

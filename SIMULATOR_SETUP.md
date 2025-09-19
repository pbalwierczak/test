# Simulator Setup

The simulator has been configured to run as a standalone container that connects to the main scootin-app container via Kafka events and HTTP API calls.

## Architecture

- **Main App**: Contains the API server, database, and Kafka (runs via `docker-compose.yml`)
- **Simulator**: Contains only the simulator app (runs via `docker-compose.simulator.yml`)
- **Communication**: Simulator connects to main app via:
  - **Kafka Events**: For publishing trip and location events
  - **HTTP API**: For scooter management and trip operations
- **Network**: Both containers run on the same Docker network (`scootin-network`)

## Usage

### 1. Start the Main Application

First, start the main application with database and Kafka:

```bash
make start-app
```

This will start:
- PostgreSQL database
- Apache Kafka with Zookeeper
- Scootin API server

### 2. Load Sample Data

Before starting the simulator, load sample data into the database:

```bash
make seed
```

This will load:
- Sample users
- Sample scooters with initial locations
- Sample trip data

> ⚠️ **Important**: The simulator relies on existing data in the database. Without sample data, the simulator won't have scooters or users to work with.

### 3. Start the Simulator

Once the main app is running and sample data is loaded, start the simulator:

```bash
make start-sim
```

The simulator will:
- Build the simulator Docker image
- Connect to the running scootin-app container via Docker network
- Start simulating scooter and user behavior using existing data
- Publish events to Kafka topics

### 4. Monitor Services

View logs and check status:

```bash
make logs-app      # View application logs
make logs-sim      # View simulator logs
make status        # Check all running services
```

## Configuration

The simulator connects to the main app using environment variables:

- `SIMULATOR_SERVER_URL`: HTTP API endpoint (default: `http://scootin-app:8080`)
- `KAFKA_BROKERS`: Kafka broker addresses (default: `kafka:29092`)

## Docker Compose Files

- `docker-compose.yml`: Main application with database, Kafka, and API server
- `docker-compose.simulator.yml`: Simulator only (connects to main app network)

## Event-Driven Communication

The simulator uses Kafka events for real-time communication:

### Kafka Topics
- `scooter.trip.started` - When a user begins a trip
- `scooter.trip.ended` - When a trip is completed  
- `scooter.location.updated` - Periodic location updates during trips

### Event Flow
1. **Simulator** → Publishes events to Kafka topics
2. **Kafka** → Stores and distributes events
3. **Server** → Consumes events and processes them through existing services

## Troubleshooting

If the simulator cannot connect to the main app:

1. Ensure the main app is running: `make status`
2. Check the API is accessible: `curl http://localhost:8080/api/v1/health`
3. Check simulator logs: `make logs-sim`
4. Verify network connectivity: `docker network ls | grep scootin`

### Common Issues

**"no such host" error**: This means the simulator can't reach the main app. Make sure:
- The main app is running (`make start-app`)
- Both containers are on the same Docker network
- The simulator is using the correct URL (`scootin-app:8080`)

**Connection refused**: If you get a connection error:
- Ensure the main app is running and listening on port 8080
- Check that the main app is accessible: `curl http://localhost:8080/api/v1/health`
- Verify Kafka is running: `docker ps | grep kafka`

**Kafka connection issues**: If events aren't being processed:
- Ensure Kafka is healthy: `docker ps | grep kafka`
- Check Kafka logs: `docker logs scootin-kafka`
- Verify topic creation: `docker exec scootin-kafka kafka-topics --bootstrap-server localhost:9092 --list`

## Available Commands

**Quick Start:**
- `make start-app` - Start app with database and Kafka
- `make seed` - Load sample data (required before simulator)
- `make start-sim` - Start simulator
- `make logs-app` - View app logs
- `make logs-sim` - View simulator logs
- `make status` - Show running services
- `make kill-all` - Stop everything

**Development:**
- `make clean` - Clean up everything (containers, images, volumes)

## Benefits

- **Separation of Concerns**: Simulator and main app are independent
- **Resource Efficiency**: Simulator doesn't need its own database or Kafka
- **Scalability**: Can run multiple simulators against one main app
- **Development**: Easier to develop and test simulator independently
- **Event-Driven**: Real-time communication via Kafka events
- **Network Isolation**: Secure communication within Docker network

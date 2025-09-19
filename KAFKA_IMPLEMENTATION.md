# Kafka Implementation for Scootin' Aboot

This document describes the Kafka implementation that enables event-driven architecture for the scooter simulation system.

## Overview

The implementation provides a dual-mode operation where the simulator can either use REST API calls or Kafka events to communicate with the server. This allows for comparison between synchronous and asynchronous communication patterns.

## Architecture

### Event Flow

1. **Simulator** → Publishes events to Kafka topics
2. **Kafka** → Stores and distributes events
3. **Server** → Consumes events and processes them through existing services

### Event Types

- **Trip Started**: `scooter.trip.started`
- **Trip Ended**: `scooter.trip.ended`  
- **Location Updated**: `scooter.location.updated`

## Configuration

### Environment Variables

```bash
# Mode Configuration
SIMULATOR_MODE=kafka  # or "rest" for REST mode

# Kafka Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_CLIENT_ID=scooter-simulator
KAFKA_SECURITY_PROTOCOL=PLAINTEXT
KAFKA_TOPIC_TRIP_STARTED=scooter.trip.started
KAFKA_TOPIC_TRIP_ENDED=scooter.trip.ended
KAFKA_TOPIC_LOCATION_UPDATED=scooter.location.updated
```

### Docker Compose

Kafka is already configured in `docker-compose.yml`:

```yaml
kafka:
  image: confluentinc/cp-kafka:7.4.0
  container_name: scootin-kafka
  environment:
    KAFKA_AUTO_CREATE_TOPICS_ENABLE: 'true'
  ports:
    - "9092:9092"
```

## Usage

### Running in Kafka Mode

1. Start the infrastructure:
   ```bash
   docker-compose up -d postgres kafka
   ```

2. Set environment variable:
   ```bash
   export SIMULATOR_MODE=kafka
   ```

3. Start the server:
   ```bash
   go run cmd/server/main.go
   ```

4. Start the simulator:
   ```bash
   go run cmd/simulator/main.go
   ```

### Running in REST Mode (Default)

1. Start the infrastructure:
   ```bash
   docker-compose up -d postgres
   ```

2. Start the server:
   ```bash
   go run cmd/server/main.go
   ```

3. Start the simulator:
   ```bash
   go run cmd/simulator/main.go
   ```

## Event Schemas

### Trip Started Event

```json
{
  "eventType": "trip.started",
  "eventId": "uuid",
  "timestamp": "2024-01-01T12:00:00Z",
  "version": "1.0",
  "data": {
    "tripId": "uuid",
    "scooterId": "uuid",
    "userId": "uuid",
    "startLatitude": 45.4215,
    "startLongitude": -75.6972,
    "startTime": "2024-01-01T12:00:00Z"
  }
}
```

### Trip Ended Event

```json
{
  "eventType": "trip.ended",
  "eventId": "uuid",
  "timestamp": "2024-01-01T12:15:00Z",
  "version": "1.0",
  "data": {
    "tripId": "uuid",
    "scooterId": "uuid",
    "userId": "uuid",
    "endLatitude": 45.4300,
    "endLongitude": -75.6800,
    "endTime": "2024-01-01T12:15:00Z",
    "durationSeconds": 900
  }
}
```

### Location Updated Event

```json
{
  "eventType": "location.updated",
  "eventId": "uuid",
  "timestamp": "2024-01-01T12:05:00Z",
  "version": "1.0",
  "data": {
    "scooterId": "uuid",
    "tripId": "uuid",
    "latitude": 45.4250,
    "longitude": -75.6900,
    "heading": 45.5,
    "speed": 15.2
  }
}
```

## Implementation Details

### Producer (Simulator Side)

- **Location**: `pkg/kafka/producer.go`
- **Interface**: `EventProducer`
- **Implementations**: `KafkaProducer`, `MockProducer`

### Consumer (Server Side)

- **Location**: `internal/kafka/consumer.go`
- **Integration**: Consumes events and calls existing service methods

### Abstract Publisher

- **Location**: `pkg/simulator/publisher.go`
- **Interface**: `EventPublisher`
- **Implementations**: `KafkaEventPublisher`, `RESTEventPublisher`

## Testing

Run the Kafka tests:

```bash
go test ./pkg/kafka/...
```

## Benefits of Kafka Mode

1. **Decoupling**: Simulator and server are decoupled through events
2. **Scalability**: Multiple consumers can process events
3. **Reliability**: Events are persisted and can be replayed
4. **Observability**: Event flow is visible and traceable
5. **Flexibility**: Easy to add new event consumers

## Monitoring

The implementation includes comprehensive logging for:
- Event publishing
- Event consumption
- Error handling
- Performance metrics

Check the logs for event flow and any issues:

```bash
# Server logs
docker logs scootin-app

# Simulator logs
go run cmd/simulator/main.go
```

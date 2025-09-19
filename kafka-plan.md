# Kafka Migration Plan - Scooter Simulator POC

## Overview

This document outlines the migration plan from REST-based scooter simulation to event-driven architecture using Apache Kafka. The goal is to replace direct HTTP calls for trip management and location updates with Kafka events, enabling better scalability, decoupling, and real-time processing.

## Current Architecture Analysis

### Current REST Endpoints Being Migrated
- `POST /api/v1/scooters/{id}/trip/start` - Start a trip
- `POST /api/v1/scooters/{id}/trip/end` - End a trip  
- `POST /api/v1/scooters/{id}/location` - Update scooter location

### Current Simulator Components
- **Simulator**: Orchestrates the simulation, manages scooters and users
- **Scooter**: Simulates individual scooter behavior (movement, trip lifecycle)
- **User**: Simulates user behavior (searching for scooters)
- **APIClient**: Handles HTTP communication with the server

## Migration Strategy

### Phase 1: Event Schema Design & Infrastructure Setup
**Goal**: Define event schemas and set up Kafka infrastructure

#### 1.1 Event Schemas
Define Avro/JSON schemas for the three main events:

```json
// Trip Start Event
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

// Trip End Event
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

// Location Update Event
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

#### 1.2 Kafka Topics
- `scooter.trip.started` - Trip start events
- `scooter.trip.ended` - Trip end events  
- `scooter.location.updated` - Location update events
- `scooter.events` - All scooter events (alternative single topic approach)

#### 1.3 Infrastructure Updates
- Add Kafka client dependencies to Go modules
- Update docker-compose to include Kafka (already present)
- Add Kafka configuration to simulator config

### Phase 2: Kafka Producer Implementation
**Goal**: Replace REST calls with Kafka event publishing

#### 2.1 Create Kafka Producer Package
```
pkg/kafka/
├── producer.go          # Main producer interface
├── events.go           # Event struct definitions
├── config.go           # Kafka configuration
└── producer_test.go    # Unit tests
```

#### 2.2 Event Producer Interface
```go
type EventProducer interface {
    PublishTripStarted(ctx context.Context, event TripStartedEvent) error
    PublishTripEnded(ctx context.Context, event TripEndedEvent) error
    PublishLocationUpdated(ctx context.Context, event LocationUpdatedEvent) error
    Close() error
}
```

#### 2.3 Update Simulator Components
- Replace `APIClient` calls in `Scooter` with Kafka producer calls
- Maintain same business logic, only change communication layer
- Add retry logic and error handling for Kafka operations

### Phase 3: Kafka Consumer Implementation (Server Side)
**Goal**: Process Kafka events on the server side

#### 3.1 Create Kafka Consumer Package
```
internal/kafka/
├── consumer.go         # Main consumer interface
├── handlers.go         # Event handlers
├── config.go          # Consumer configuration
└── consumer_test.go   # Unit tests
```

#### 3.2 Event Handlers
- `handleTripStarted()` - Process trip start events
- `handleTripEnded()` - Process trip end events
- `handleLocationUpdated()` - Process location update events

#### 3.3 Integration with Existing Services
- Wire Kafka consumers to existing `TripService` and `ScooterService`
- Maintain existing business logic and validation
- Add event processing metrics and monitoring

### Phase 4: Dual-Mode Operation (POC Phase)
**Goal**: Support both REST and Kafka modes for comparison

#### 4.1 Configuration-Based Switching
```go
type Config struct {
    // ... existing config
    SimulatorMode string // "rest" or "kafka"
    KafkaConfig   KafkaConfig
}
```

#### 4.2 Abstract Communication Layer
```go
type EventPublisher interface {
    PublishTripStarted(ctx context.Context, event TripStartedEvent) error
    PublishTripEnded(ctx context.Context, event TripEndedEvent) error
    PublishLocationUpdated(ctx context.Context, event LocationUpdatedEvent) error
}

type RESTPublisher struct { /* ... */ }
type KafkaPublisher struct { /* ... */ }
```

### Phase 5: Monitoring & Observability
**Goal**: Add comprehensive monitoring for event-driven architecture

#### 5.1 Logging
- Structured logging for all events
- Correlation IDs for tracing across services
- Error tracking and alerting

#### 5.2 Health Checks
- Kafka connectivity health checks
- Consumer group health monitoring
- Event processing pipeline health

## Implementation Details

### Dependencies
```go
// Add to go.mod
require (
    github.com/Shopify/sarama v1.38.1
    github.com/IBM/sarama v1.41.2
    github.com/confluentinc/confluent-kafka-go v2.2.0+incompatible
)
```

### Configuration Structure
```go
type KafkaConfig struct {
    Brokers          []string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
    ClientID         string   `env:"KAFKA_CLIENT_ID" envDefault:"scooter-simulator"`
    SecurityProtocol string   `env:"KAFKA_SECURITY_PROTOCOL" envDefault:"PLAINTEXT"`
    Topics           KafkaTopics
}

type KafkaTopics struct {
    TripStarted     string `env:"KAFKA_TOPIC_TRIP_STARTED" envDefault:"scooter.trip.started"`
    TripEnded       string `env:"KAFKA_TOPIC_TRIP_ENDED" envDefault:"scooter.trip.ended"`
    LocationUpdated string `env:"KAFKA_TOPIC_LOCATION_UPDATED" envDefault:"scooter.location.updated"`
}
```

### Error Handling Strategy
1. **Retry Logic**: Exponential backoff for transient failures
2. **Dead Letter Queue**: Failed events after max retries
3. **Circuit Breaker**: Prevent cascade failures

### Testing Strategy
1. **Unit Tests**: Individual component testing

### Data Consistency
- Kafka events are idempotent
- Database state remains consistent
- No data loss during rollback
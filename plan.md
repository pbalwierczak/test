# Scootin' Aboot - Electric Scooter Management System

Important: see plan-support.md for additional details.
When implementing a point, mark it as done with x in [ ].

## Project Overview
A backend service for managing electric scooters in Ottawa and Montreal, providing REST API for scooter event collection and mobile client reporting.
These are the only requirements - no need to invent more:
- The scooters report an event when a trip begins, report an event when the trip ends, and send in periodic updates on their location.
- After beginning a trip, the scooter is considered occupied. 
- After a trip ends the scooter becomes free for use. 
- A location update must contain the time, and geographical coordinates.
- Mobile clients can query scooter locations and statuses in any rectangular location (e.g. two pair of coordinates), and filter them by status.
- While there will be no actual mobile clients, implement child process that would start with main process and spawn three fake clients using API randomly (finding scooters, travelling for 10-15 seconds whilst updating location every 3 seconds, and resting for 2-5 seconds before starting next trip). Client movement in straight line will be good enough.
- Both scooters and mobile client users can be identified by an UUID.
- For the sake of simplicity, both mobile client apps and scooters can authenticate with the server using a static API key (i.e. no individual credentials necessary).

## Architecture & Technology Stack

### Core Technologies
- **Language**: Go 1.21+
- **Web Framework**: Gin (lightweight, fast HTTP router)
- **Database**: PostgreSQL (ACID compliance, JSON support for coordinates)
- **ORM**: GORM (Go ORM with migration support)
- **Authentication**: Static API key (no individual credentials)
- **Containerization**: Docker & Docker Compose
- **Testing**: Go testing package + testify
- **Documentation**: Swagger/OpenAPI 3.0

### Project Structure
```
scootin-aboot/
├── cmd/
│   ├── server/
│   │   └── main.go
│   └── simulator/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── routes/
│   ├── models/
│   ├── services/
│   ├── repository/
│   └── config/
├── pkg/
│   ├── auth/
│   │   └── apikey/
│   ├── database/
│   ├── simulator/
│   └── utils/
├── migrations/
│   ├── 001_create_scooters_table.up.sql
│   ├── 001_create_scooters_table.down.sql
│   ├── 002_create_trips_table.up.sql
│   ├── 002_create_trips_table.down.sql
│   ├── 003_create_location_updates_table.up.sql
│   ├── 003_create_location_updates_table.down.sql
│   ├── 004_create_users_table.up.sql
│   └── 004_create_users_table.down.sql
├── seeds/
│   ├── scooters.sql
│   ├── users.sql
│   └── sample_trips.sql
├── docker/
├── docs/
├── scripts/
├── .env
├── .env.example
├── Makefile
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── README.md
└── plan.md
```

## Phase 1: Project Setup & Foundation (Session 1)

### 1.1 Project Initialization ✅
- [x] Initialize git repo and Go module (`go mod init scootin-aboot`)
- [x] Set up project directory structure
- [x] Create `.env` and `.env.example` files for configuration
- [x] Set up environment variable loading (godotenv or viper)
- [x] Set up logging (zap)
- [x] Add godotenv dependency for .env file loading

### 1.2 Database Design
- [x] Design database schema for:
  - Scooters table (id, status, current_location)
  - Trips table (id, scooter_id, user_id, start_time, end_time, start_location, end_location)
  - Location_updates table (id, trip_id, latitude, longitude, timestamp)
  - Users table (id, created_at)
- [x] Create migration files using golang-migrate
- [x] Set up database connection and configuration
- [x] Create seed data files and seed data load scripts for development

### 1.3 Basic API Structure
- [x] Create basic middleware (logging, recovery)
- [x] Set up basic error handling

### 1.4 Testing Foundation
- [ ] Set up testing framework (Go testing package + testify)
- [ ] Add test configuration and setup
- [ ] Test basic configuration loading
- [ ] Test health check endpoint
- [ ] Add running tests to makefile

## Phase 2: Dockerization & Container Setup (Session 2)

### 2.1 Docker Configuration
- [ ] Create Dockerfile for main application
- [ ] Create Dockerfile for simulator
- [ ] Set up docker-compose.yml with both services
- [ ] Add environment configuration for Docker
- [ ] Configure database service in docker-compose
- [ ] Set up volume mounts for development

### 2.2 Container Management
- [ ] Add Docker targets to Makefile
- [ ] Create development environment setup
- [ ] Test container builds and startup
- [ ] Ensure database migrations work in containers
- [ ] Verify seed data loading in containers

## Phase 3: Core API Implementation (Session 3)

### 3.1 Authentication & Authorization
- [ ] Implement static API key validation (single key for entire system)
- [ ] Create API key middleware for request authentication
- [ ] Load API key from `.env` file using environment variable loader

### 3.2 Scooter Management API
- [ ] `POST /api/v1/scooters/{id}/trip/start` - Start trip
- [ ] `POST /api/v1/scooters/{id}/trip/end` - End trip
- [ ] `POST /api/v1/scooters/{id}/location` - Update location
- [ ] `GET /api/v1/scooters` - Query scooters with filters
- [ ] `GET /api/v1/scooters/{id}` - Get specific scooter details
- [ ] `GET /api/v1/scooters/closest` - Find closest scooters
- [ ] Secure endpoints with API key middleware

### 3.3 Data Models & Repository Layer
- [ ] Implement GORM models
- [ ] Create repository interfaces and implementations
- [ ] Add database migrations
- [ ] Implement basic CRUD operations

### 3.4 API Testing
- [ ] Test authentication middleware
- [ ] Test all scooter management endpoints
- [ ] Test API key validation
- [ ] Test error handling and responses
- [ ] Test repository layer CRUD operations
- [ ] Test service layer business logic
- [ ] Test geographic search functionality

## Phase 4: Business Logic & Services (Session 4)

### 4.1 Trip Management Service
- [ ] Implement trip lifecycle management
- [ ] Add scooter status validation
- [ ] Handle concurrent trip attempts
- [ ] Implement location update validation

### 4.2 Query & Filtering Service
- [ ] Implement rectangular area filtering
- [ ] Add status-based filtering
- [ ] Implement closest scooter search with distance calculation
- [ ] Add geographic indexing for performance
- [ ] Optimize database queries with indexes
- [ ] Add pagination support

### 4.3 Location Services
- [ ] Validate geographical coordinates
- [ ] Implement distance calculations (Haversine formula)
- [ ] Add location history tracking
- [ ] Implement geographic bounding box queries
- [ ] Add PostGIS integration for advanced geographic queries

### 4.4 Service Layer Testing
- [ ] Test trip lifecycle management
- [ ] Test scooter status validation logic
- [ ] Test concurrent trip handling
- [ ] Test location update validation
- [ ] Test geographic filtering algorithms
- [ ] Test distance calculation accuracy
- [ ] Test database query optimization
- [ ] Test pagination functionality

## Phase 5: Advanced Simulator Implementation (Session 5)

### 5.1 Simulator Architecture
- [ ] Create separate `cmd/simulator` program
- [ ] Implement concurrent user and scooter simulation
- [ ] Add configurable simulation parameters (20 scooters, 5 users by default)
- [ ] Implement graceful shutdown with Ctrl+C handling

### 5.2 User Simulation
- [ ] **User Behavior**: Find available scooter → Start trip → Drive around → End trip → Rest
- [ ] **Trip Duration**: 5-10 seconds (configurable)
- [ ] **Driving Speed**: 100 km/h (27.78 m/s) for realistic movement
- [ ] **Movement Pattern**: Straight-line movement only (no direction changes)
- [ ] **Rest Period**: 2-5 seconds between trips

### 5.3 Scooter Simulation
- [ ] **Location Updates**: Every 3 seconds during trips
- [ ] **Status Management**: Available → Occupied → Available
- [ ] **Movement Calculation**: Realistic GPS coordinate updates based on speed
- [ ] **Geographic Bounds**: Ottawa and Montreal area coordinates

### 5.4 Movement Physics
- [ ] **Speed Calculation**: 100 km/h = 27.78 m/s
- [ ] **Distance per Update**: ~83.33 meters per 3-second update
- [ ] **Coordinate Conversion**: Meters to GPS degrees (latitude/longitude)
- [ ] **Direction Changes**: No direction changes - straight-line movement only

### 5.5 Process Management
- [ ] **Concurrent Goroutines**: Separate goroutines for each user and scooter
- [ ] **Signal Handling**: Graceful shutdown on SIGINT/SIGTERM
- [ ] **Status Monitoring**: Real-time simulation statistics
- [ ] **API Integration**: HTTP client for server communication

### 5.6 Simulator Testing
- [ ] Test simulator architecture and initialization
- [ ] Test user behavior simulation
- [ ] Test scooter movement physics
- [ ] Test concurrent goroutine management
- [ ] Test signal handling and graceful shutdown
- [ ] Test API integration and error handling
- [ ] Test configuration parameter validation
- [ ] Test movement calculation accuracy

## Phase 6: Testing & Quality Assurance (Session 6)

### 6.1 Test Coverage & Quality
- [ ] Achieve >80% code coverage across all packages
- [ ] Review and enhance existing unit tests
- [ ] Add missing test cases for edge scenarios
- [ ] Optimize test execution performance

### 6.2 Integration Testing
- [ ] Test complete API workflows end-to-end
- [ ] Test database operations with real data
- [ ] Test concurrent operations and race conditions
- [ ] Test error scenarios and recovery

### 6.3 Load Testing
- [ ] Test with multiple concurrent clients
- [ ] Test database performance under load
- [ ] Optimize slow queries identified during testing
- [ ] Test simulator under high load conditions

### 6.4 Test Infrastructure
- [ ] Set up test database for integration tests
- [ ] Create test data fixtures and factories
- [ ] Implement test cleanup and teardown
- [ ] Add test reporting and coverage analysis

## Phase 7: Makefile & Build System (Session 7)

### 7.1 Makefile Implementation
- [ ] Create comprehensive Makefile with targets:
  - `make server` - Run the main API server
  - `make simulator` - Run the simulation program
  - `make build` - Build both server and simulator
  - `make test` - Run all tests (unit + integration)
  - `make test-unit` - Run unit tests only
  - `make test-integration` - Run integration tests only
  - `make test-coverage` - Run tests with coverage report
  - `make docker` - Build Docker images
  - `make clean` - Clean build artifacts
  - `make migrate-up` - Run database migrations
  - `make migrate-down` - Rollback database migrations
  - `make migrate-reset` - Reset database (down + up)
  - `make seed` - Insert seed data
  - `make setup` - Full setup (migrate + seed)
  - `make dev` - Start development environment (server + simulator)
- [ ] Add configuration options for simulator parameters
- [ ] Implement graceful shutdown handling

### 7.2 Simulator Configuration
- [ ] **Command Line Flags** (override .env defaults):
  - `--scooters` - Number of scooters (default from .env)
  - `--users` - Number of users (default from .env)
  - `--server-url` - API server URL (default from .env)
  - `--speed` - Driving speed in km/h (default from .env)
  - `--trip-duration` - Trip duration range (default from .env)
- [ ] **Environment Variables**: Load all defaults from `.env` file
- [ ] **Configuration Priority**: CLI flags > .env file > hardcoded defaults

### 7.3 Database Management
- [ ] **Migration System**: Use golang-migrate or similar for database migrations
- [ ] **Migration Files**: Create SQL migration files in `migrations/` directory
- [ ] **Seed Data**: Create comprehensive seed data for development and testing
- [ ] **Database Connection**: Ensure migrations work with both local and Docker databases

### 7.4 Seed Data Implementation
- [ ] **Scooter Seed Data**: Pre-populate with 20+ scooters in Ottawa and Montreal
  - Realistic GPS coordinates within city boundaries
  - Mix of statuses: available, occupied
  - Varied last_seen timestamps
- [ ] **User Seed Data**: Create test users (no API key storage needed)
  - 5+ test users for trip tracking
  - Single static API key for entire system (stored in environment)
  - Different user types (regular, admin, test)
- [ ] **Geographic Distribution**: Spread scooters across realistic locations
  - Ottawa: Parliament Hill, ByWard Market, Rideau Centre areas
  - Montreal: Old Port, Downtown, Plateau areas
  - 10km radius from city centers
- [ ] **Status Variety**: Mix of available, occupied scooters
- [ ] **Historical Data**: Sample trips and location updates for testing
  - Completed trips with full location history
  - Ongoing trips for testing
  - Various trip durations and distances

### 7.5 Process Management
- [ ] **Signal Handling**: Proper SIGINT/SIGTERM handling
- [ ] **Resource Cleanup**: Graceful shutdown of all goroutines
- [ ] **Status Reporting**: Real-time statistics display
- [ ] **Error Recovery**: Handle API failures gracefully

## Phase 8: Documentation & Deployment (Session 8)

### 8.1 API Documentation
- [ ] Generate Swagger/OpenAPI documentation
- [ ] Document all endpoints with examples
- [ ] Create API usage guide

### 8.2 Documentation
- [ ] Complete README.md with setup instructions
- [ ] Create ASSUMPTIONS.md
- [ ] Add database schema documentation
- [ ] Create deployment guide
- [ ] Document simulator usage and configuration
- [ ] **Maintain ASSUMPTIONS.md** - Update as new assumptions are made or existing ones are validated/invalidated

## Phase 9: Final Polish & Optimization (Session 9)

### 9.1 Performance Optimization
- [ ] Database query optimization
- [ ] Add caching where appropriate
- [ ] Optimize API response times
- [ ] Memory usage optimization
- [ ] Simulator performance tuning

### 9.2 Monitoring & Observability
- [ ] Add structured logging
- [ ] Implement metrics collection
- [ ] Add health check endpoints
- [ ] Create monitoring dashboard
- [ ] Simulator statistics and metrics

### 9.3 Security & Validation
- [ ] Input validation and sanitization
- [ ] Rate limiting
- [ ] SQL injection prevention
- [ ] Security headers
- [ ] API key validation (static key from environment)

---

*This plan will be updated as the project progresses and requirements evolve.*

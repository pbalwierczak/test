## API Endpoints Specification

### Scooter Management
```
POST   /api/v1/scooters/{id}/trip/start    # Requires API key
POST   /api/v1/scooters/{id}/trip/end      # Requires API key
POST   /api/v1/scooters/{id}/location      # Requires API key
GET    /api/v1/scooters?lat_min=45.0&lat_max=46.0&lng_min=-76.0&lng_max=-75.0&status=available  # Requires API key
GET    /api/v1/scooters/{id}               # Requires API key
GET    /api/v1/scooters/closest?lat=45.4215&lng=-75.6972&radius=1000&status=available&limit=10  # Requires API key
```

### System
```
GET    /api/v1/health                      # Public endpoint
GET    /api/v1/metrics                     # Requires API key
```

### Authentication
- **Header**: `X-API-Key: your-static-api-key`
- **Static Key**: Single key for entire system (stored in environment variables)
- **No Database Storage**: API key not stored in database, only in configuration
- **No Individual Credentials**: Simplified authentication as per requirements

## Database Schema

### Scooters Table
```sql
CREATE TABLE scooters (
    id UUID PRIMARY KEY,
    status VARCHAR(20) NOT NULL, -- 'available', 'occupied', 'maintenance'
    current_latitude DECIMAL(10, 8),
    current_longitude DECIMAL(11, 8),
    last_seen TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Trips Table
```sql
CREATE TABLE trips (
    id UUID PRIMARY KEY,
    scooter_id UUID REFERENCES scooters(id),
    user_id UUID NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    start_latitude DECIMAL(10, 8),
    start_longitude DECIMAL(11, 8),
    end_latitude DECIMAL(10, 8),
    end_longitude DECIMAL(11, 8),
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Location Updates Table
```sql
CREATE TABLE location_updates (
    id UUID PRIMARY KEY,
    trip_id UUID REFERENCES trips(id),
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    user_type VARCHAR(50) DEFAULT 'regular', -- 'regular', 'admin', 'test'
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Authentication
- **Static API Key**: Single key for entire system (stored in `.env` file)
- **No Database Storage**: API key not stored in database, only in configuration
- **Header**: `X-API-Key: your-static-api-key`
- **Environment File**: `.env` file for all configuration including API key

## Scooter Search Implementation

### Closest Scooter Search Algorithm
The system will implement multiple search strategies for finding scooters:

#### 1. Geographic Bounding Box Search (Fast)
```sql
-- Find scooters within rectangular area
SELECT * FROM scooters 
WHERE latitude BETWEEN ? AND ? 
  AND longitude BETWEEN ? AND ?
  AND status = 'available'
ORDER BY 
  (latitude - ?) * (latitude - ?) + (longitude - ?) * (longitude - ?)
LIMIT ?;
```

#### 2. Haversine Distance Calculation (Accurate)
```go
// Haversine formula for accurate distance calculation
func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
    const R = 6371000 // Earth's radius in meters
    dLat := (lat2 - lat1) * math.Pi / 180
    dLng := (lng2 - lng1) * math.Pi / 180
    a := math.Sin(dLat/2)*math.Sin(dLat/2) + 
         math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180) * 
         math.Sin(dLng/2)*math.Sin(dLng/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    return R * c
}
```

#### 3. Database Indexing Strategy
```sql
-- Geographic index for fast spatial queries
CREATE INDEX idx_scooters_location ON scooters USING GIST (
    ST_Point(longitude, latitude)
);

-- Composite index for status + location
CREATE INDEX idx_scooters_status_location ON scooters (status, latitude, longitude);
```

### Search API Endpoints

#### Rectangular Area Search
```
GET /api/v1/scooters?lat_min=45.0&lat_max=46.0&lng_min=-76.0&lng_max=-75.0&status=available
```
- **Use Case**: Map view with rectangular bounds
- **Performance**: Very fast with proper indexing
- **Accuracy**: Good for large areas

#### Closest Scooter Search
```
GET /api/v1/scooters/closest?lat=45.4215&lng=-75.6972&radius=1000&status=available&limit=10
```
- **Parameters**:
  - `lat`, `lng`: User's current location
  - `radius`: Search radius in meters (default: 1000m)
  - `status`: Filter by scooter status (optional)
  - `limit`: Maximum results (default: 10)
- **Use Case**: "Find nearest available scooter"
- **Performance**: Optimized with geographic indexing
- **Accuracy**: Precise distance calculation

### Search Performance Optimization

#### 1. Two-Phase Search Strategy
```go
// Phase 1: Fast bounding box search
func FindScootersInBounds(minLat, maxLat, minLng, maxLng float64) []Scooter {
    // Use spatial index for fast rectangular search
}

// Phase 2: Accurate distance calculation
func SortByDistance(scooters []Scooter, userLat, userLng float64) []Scooter {
    // Calculate precise distances and sort
}
```

#### 2. Caching Strategy
- **Redis Cache**: Cache frequent search results
- **TTL**: 30 seconds for location-based searches
- **Key Pattern**: `scooters:search:{lat}:{lng}:{radius}:{status}`

#### 3. Database Optimization
- **PostGIS Extension**: For advanced geographic queries
- **Spatial Indexes**: GIST indexes for geographic data
- **Query Optimization**: Use EXPLAIN ANALYZE for query tuning

### Search Response Format
```json
{
  "scooters": [
    {
      "id": "uuid",
      "latitude": 45.4215,
      "longitude": -75.6972,
      "status": "available",
      "distance_meters": 150,
      "last_seen": "2025-01-27T10:30:00Z"
    }
  ],
  "total": 5,
  "search_center": {
    "latitude": 45.4215,
    "longitude": -75.6972
  },
  "search_radius": 1000
}
```

## Configuration Management

### .env File Structure
```bash
# API Configuration
API_KEY=test-api-key-12345
SERVER_PORT=8080
SERVER_HOST=localhost

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=scootin_aboot
DB_USER=postgres
DB_PASSWORD=password
DB_SSLMODE=disable

# Simulator Configuration
SIMULATOR_SCOOTERS=20
SIMULATOR_USERS=5
SIMULATOR_SERVER_URL=http://localhost:8080
SIMULATOR_SPEED=100
SIMULATOR_TRIP_DURATION_MIN=5
SIMULATOR_TRIP_DURATION_MAX=10
SIMULATOR_REST_MIN=2
SIMULATOR_REST_MAX=5

# Geographic Configuration (now constants in code)
# Ottawa: 45.4215°N, 75.6972°W
# Montreal: 45.5017°N, 73.5673°W
# City Radius: 10km

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### .env.example File
- [ ] Create `.env.example` with all configuration variables
- [ ] Document each variable with comments
- [ ] Provide sensible defaults for development
- [ ] Exclude sensitive values (use placeholder values)

## Success Criteria

### Functional Requirements
- [ ] Scooters can start and end trips
- [ ] Location updates are recorded during trips
- [ ] Mobile clients can query scooters by location and status
- [ ] Advanced simulator simulates realistic usage patterns
- [ ] System handles concurrent operations
- [ ] Simulator runs independently with configurable parameters
- [ ] Graceful shutdown handling for both server and simulator

### Non-Functional Requirements
- [ ] API response time < 200ms for 95% of requests
- [ ] System supports 100+ concurrent users
- [ ] Database queries are optimized
- [ ] Code coverage > 80%
- [ ] Docker setup works out of the box

### Documentation Requirements
- [ ] Complete README with setup instructions
- [ ] API documentation with examples
- [ ] Database schema documentation
- [ ] Assumptions and design decisions documented

## Risk Mitigation

### Technical Risks
- **Database Performance**: Use proper indexing and query optimization
- **Concurrency Issues**: Implement proper locking and transaction management

### Project Risks
- **Scope Creep**: Stick to MVP features
- **Time Constraints**: Prioritize core functionality over nice-to-have features
- **Integration Issues**: Test early and often with Docker environment

## Simulator Technical Specifications

### Movement Physics
- **Update Interval**: 3 seconds
- **Coordinate Precision**: 6 decimal places (GPS accuracy)
- **Geographic Bounds**: 
  - Ottawa: 45.4215°N, 75.6972°W (center)
  - Montreal: 45.5017°N, 73.5673°W (center)
  - Radius: ~10km from city centers

### Simulator Behavior
- **Default Configuration**: 100 scooters, 25 users
- **Trip Duration**: 5-10 seconds (random)
- **Rest Period**: 2-5 seconds between trips
- **Direction Changes**: No direction changes - straight-line movement only
- **API Calls**: Every 3 seconds during trips
- **Concurrent Operations**: All users and scooters run simultaneously

### Makefile Targets
Should prioritize docker, all functions that are not run via docker should start with underscore _

## Development Guidelines

### Code Standards
- Follow Go best practices and idioms
- Use meaningful variable and function names
- Add comprehensive comments for complex logic
- Implement proper error handling
- Write tests for all public functions

### Git Workflow
- Use feature branches for new functionality
- Write descriptive commit messages
- Create pull requests for code review
- Tag releases with semantic versioning

### Testing Strategy

#### Test Organization
- **Unit Tests**: Each package has corresponding `*_test.go` files
No other testing is necessary

#### Test Categories
- **Unit Tests**: Test individual functions and methods in isolation, mock what is needed
No need for other types of tests

#### Test Coverage Requirements
- **Minimum Coverage**: 80% code coverage across all packages
- **Critical Paths**: 100% coverage for business logic and API handlers
- **Edge Cases**: Comprehensive testing of error scenarios and boundary conditions

#### Test Execution
- **Local Development**: `make _test` runs all tests locally
- **Docker Environment**: `make test` runs tests in Docker container
- **Coverage Reports**: `make test-coverage` generates detailed coverage reports
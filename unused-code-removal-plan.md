# Unused Code Removal Plan

## Overview
This document outlines the unused methods and functions identified in the codebase that can be safely removed. The analysis was performed by searching for actual usage patterns across the entire codebase, excluding unit test usage.

## Analysis Methodology
- Searched for method calls across all Go files
- Excluded test files from usage analysis
- Verified that methods are not used in production code
- Checked routing configuration for handler methods
- Analyzed service layer method usage

## Repository Layer - Unused Methods

### TripRepository Interface
**File**: `internal/repository/trip.go`

**Unused Methods** (can be removed):
- `GetByStatus(ctx context.Context, status models.TripStatus) ([]*models.Trip, error)`
- `GetActive(ctx context.Context) ([]*models.Trip, error)`
- `GetCompleted(ctx context.Context) ([]*models.Trip, error)`
- `GetCancelled(ctx context.Context) ([]*models.Trip, error)`
- `GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Trip, error)`
- `GetByScooterID(ctx context.Context, scooterID uuid.UUID) ([]*models.Trip, error)`
- `GetByDateRange(ctx context.Context, start, end time.Time) ([]*models.Trip, error)`
- `GetByUserIDAndDateRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*models.Trip, error)`
- `GetByScooterIDAndDateRange(ctx context.Context, scooterID uuid.UUID, start, end time.Time) ([]*models.Trip, error)`
- `GetTripCount(ctx context.Context) (int64, error)`
- `GetTripCountByUser(ctx context.Context, userID uuid.UUID) (int64, error)`
- `GetTripCountByScooter(ctx context.Context, scooterID uuid.UUID) (int64, error)`

### ScooterRepository Interface
**File**: `internal/repository/scooter.go`

**Unused Methods** (can be removed):
- `GetAvailable(ctx context.Context) ([]*models.Scooter, error)`
- `GetOccupied(ctx context.Context) ([]*models.Scooter, error)`
- `GetInRadius(ctx context.Context, latitude, longitude, radiusKm float64) ([]*models.Scooter, error)`
- `GetAvailableInBounds(ctx context.Context, minLat, maxLat, minLng, maxLng float64) ([]*models.Scooter, error)`

### LocationUpdateRepository Interface
**File**: `internal/repository/location_update.go`

**Unused Methods** (can be removed):
- `GetByScooterIDOrdered(ctx context.Context, scooterID uuid.UUID) ([]*models.LocationUpdate, error)`
- `GetLatestByScooterID(ctx context.Context, scooterID uuid.UUID) (*models.LocationUpdate, error)`
- `GetByDateRange(ctx context.Context, start, end time.Time) ([]*models.LocationUpdate, error)`
- `GetByScooterIDAndDateRange(ctx context.Context, scooterID uuid.UUID, start, end time.Time) ([]*models.LocationUpdate, error)`
- `GetInBounds(ctx context.Context, minLat, maxLat, minLng, maxLng float64) ([]*models.LocationUpdate, error)`
- `GetInRadius(ctx context.Context, latitude, longitude, radiusKm float64) ([]*models.LocationUpdate, error)`
- `GetUpdateCount(ctx context.Context) (int64, error)`
- `GetUpdateCountByScooter(ctx context.Context, scooterID uuid.UUID) (int64, error)`

### UserRepository Interface
**File**: `internal/repository/user.go`

**All Methods Unused** (can be removed):
- `Create(ctx context.Context, user *models.User) error`
- `GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)`
- `Update(ctx context.Context, user *models.User) error`
- `Delete(ctx context.Context, id uuid.UUID) error`
- `List(ctx context.Context, limit, offset int) ([]*models.User, error)`

**Note**: User management is not implemented in the current codebase.

## Service Layer - All Methods Used
**Files**: `internal/services/trip_service.go`, `internal/services/scooter_service.go`

All service interface methods are being used in handlers or other services. No removal needed.

## Handler Layer - All Methods Used
**Files**: `internal/api/handlers/*.go`

All handler methods are properly routed in `internal/api/routes/router.go`. No removal needed.

## Utility Functions - Unused Methods

### Logger Utilities
**File**: `pkg/utils/logger.go`

**Unused Functions** (can be removed):
- `InitLogger(level, format string) error`
- `GetLogger() *zap.Logger`
- `Sync()`
- `Info(msg string, fields ...zap.Field)`
- `Debug(msg string, fields ...zap.Field)`
- `Warn(msg string, fields ...zap.Field)`
- `Error(msg string, fields ...zap.Field)`
- `Fatal(msg string, fields ...zap.Field)`

**Note**: The global `Logger` variable is also unused.

### Validation Utilities
**File**: `pkg/validation/coordinates.go`

**Unused Functions** (can be removed):
- `ValidateLatitude(lat float64) error`
- `ValidateLongitude(lng float64) error`

**Used Functions** (keep):
- `ValidateCoordinates(lat, lng float64) error`

### Database Utilities
**File**: `pkg/database/seed.go`

**Unused Functions** (can be removed):
- `SeedData(db *sql.DB, seedsPath string) error`
- `GetSeedsPath() (string, error)`

**File**: `pkg/database/connection.go`

**Unused Functions** (can be removed):
- `StartHealthCheck(db *gorm.DB, interval time.Duration) func()`

**Used Functions** (keep):
- `ConnectDatabase(dsn string) (*gorm.DB, error)`
- `AutoMigrate(db *gorm.DB) error`

### Auth Utilities
**File**: `pkg/auth/apikey/validator.go`

**All Functions Used** (keep):
- `NewValidator(apiKey string) *Validator`
- `ValidateAPIKey(providedKey string) error`
- `ExtractAPIKey(authHeader string) string`

## Implementation Plan

### Phase 1: Repository Interface Cleanup
1. Remove unused methods from repository interfaces
2. Remove corresponding implementations in GORM repositories
3. Remove corresponding mock implementations
4. Update any tests that reference removed methods

### Phase 2: Utility Function Cleanup
1. Remove unused logger utility functions
2. Remove unused validation functions
3. Remove unused database utility functions
4. Remove unused seed functions

### Phase 3: User Repository Removal
1. Remove entire UserRepository interface
2. Remove GORM implementation
3. Remove mock implementation
4. Remove user-related service dependencies
5. Update service constructors

### Phase 4: Testing and Verification
1. Run all tests to ensure nothing breaks
2. Verify that the application still compiles
3. Check that all remaining functionality works

## Files to Modify

### Repository Interfaces
- `internal/repository/trip.go`
- `internal/repository/scooter.go`
- `internal/repository/location_update.go`
- `internal/repository/user.go` (remove entirely)

### Repository Implementations
- `internal/repository/trip_gorm.go`
- `internal/repository/scooter_gorm.go`
- `internal/repository/location_update_gorm.go`
- `internal/repository/user_gorm.go` (remove entirely)

### Mock Implementations
- `internal/repository/mocks/trip_repository.go`
- `internal/repository/mocks/scooter_repository.go`
- `internal/repository/mocks/location_update_repository.go`
- `internal/repository/mocks/user_repository.go` (remove entirely)

### Service Files
- `internal/services/scooter_service.go`

## Estimated Impact
- **Lines of code to remove**: ~500-800 lines
- **Files to modify**: ~15 files
- **Files to delete**: ~4 files
- **Risk level**: Low (all removed code is confirmed unused)

## Notes
- All removals are based on actual usage analysis, not assumptions
- No production functionality will be affected
- Test coverage may decrease, but only for unused code
- Consider this cleanup as technical debt reduction

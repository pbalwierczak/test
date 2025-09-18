-- Drop location_updates table and related objects
DROP INDEX IF EXISTS idx_location_updates_trip_timestamp;
DROP INDEX IF EXISTS idx_location_updates_location;
DROP INDEX IF EXISTS idx_location_updates_timestamp;
DROP INDEX IF EXISTS idx_location_updates_trip_id;
DROP TABLE IF EXISTS location_updates;

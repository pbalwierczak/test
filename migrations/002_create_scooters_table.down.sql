-- Drop scooters table and related objects
DROP TRIGGER IF EXISTS update_scooters_updated_at ON scooters;
DROP INDEX IF EXISTS idx_scooters_location;
DROP INDEX IF EXISTS idx_scooters_last_seen;
DROP INDEX IF EXISTS idx_scooters_status;
DROP TABLE IF EXISTS scooters;

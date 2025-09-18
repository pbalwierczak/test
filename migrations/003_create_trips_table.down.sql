-- Drop trips table and related objects
DROP TRIGGER IF EXISTS update_trips_updated_at ON trips;
DROP INDEX IF EXISTS idx_trips_user_status;
DROP INDEX IF EXISTS idx_trips_scooter_status;
DROP INDEX IF EXISTS idx_trips_start_time;
DROP INDEX IF EXISTS idx_trips_status;
DROP INDEX IF EXISTS idx_trips_user_id;
DROP INDEX IF EXISTS idx_trips_scooter_id;
DROP TABLE IF EXISTS trips;

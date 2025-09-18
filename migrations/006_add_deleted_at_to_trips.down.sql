-- Remove deleted_at column from trips table
DROP INDEX IF EXISTS idx_trips_deleted_at;
ALTER TABLE trips DROP COLUMN IF EXISTS deleted_at;

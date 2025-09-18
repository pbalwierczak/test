-- Remove deleted_at column from location_updates table
DROP INDEX IF EXISTS idx_location_updates_deleted_at;
ALTER TABLE location_updates DROP COLUMN IF EXISTS deleted_at;

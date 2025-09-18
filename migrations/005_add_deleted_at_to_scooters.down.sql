-- Remove deleted_at column from scooters table
DROP INDEX IF EXISTS idx_scooters_deleted_at;
ALTER TABLE scooters DROP COLUMN IF EXISTS deleted_at;

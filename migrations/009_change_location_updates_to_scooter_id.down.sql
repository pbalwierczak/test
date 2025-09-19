-- Revert location_updates table back to trip_id reference
-- First, drop the scooter_id foreign key constraint and indexes
ALTER TABLE location_updates DROP CONSTRAINT IF EXISTS location_updates_scooter_id_fkey;
DROP INDEX IF EXISTS idx_location_updates_scooter_id;
DROP INDEX IF EXISTS idx_location_updates_scooter_timestamp;

-- Add the trip_id column back
ALTER TABLE location_updates ADD COLUMN trip_id UUID;

-- Note: We cannot easily populate trip_id from scooter_id without additional logic
-- This would require complex logic to determine which trip the location update belongs to
-- For now, we'll leave trip_id as nullable in the down migration

-- Add foreign key constraint to trips table
ALTER TABLE location_updates ADD CONSTRAINT location_updates_trip_id_fkey 
    FOREIGN KEY (trip_id) REFERENCES trips(id) ON DELETE CASCADE;

-- Create indexes for trip_id
CREATE INDEX idx_location_updates_trip_id ON location_updates(trip_id);
CREATE INDEX idx_location_updates_trip_timestamp ON location_updates(trip_id, timestamp);

-- Drop the scooter_id column
ALTER TABLE location_updates DROP COLUMN scooter_id;

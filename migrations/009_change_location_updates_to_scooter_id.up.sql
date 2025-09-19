-- Change location_updates table to reference scooters instead of trips
-- First, drop the existing foreign key constraint and index
ALTER TABLE location_updates DROP CONSTRAINT IF EXISTS location_updates_trip_id_fkey;
DROP INDEX IF EXISTS idx_location_updates_trip_id;
DROP INDEX IF EXISTS idx_location_updates_trip_timestamp;

-- Add the new scooter_id column
ALTER TABLE location_updates ADD COLUMN scooter_id UUID;

-- Populate scooter_id by joining with trips table
UPDATE location_updates 
SET scooter_id = t.scooter_id 
FROM trips t 
WHERE location_updates.trip_id = t.id;

-- Make scooter_id NOT NULL after populating
ALTER TABLE location_updates ALTER COLUMN scooter_id SET NOT NULL;

-- Add foreign key constraint to scooters table
ALTER TABLE location_updates ADD CONSTRAINT location_updates_scooter_id_fkey 
    FOREIGN KEY (scooter_id) REFERENCES scooters(id) ON DELETE CASCADE;

-- Create new indexes for scooter_id
CREATE INDEX idx_location_updates_scooter_id ON location_updates(scooter_id);
CREATE INDEX idx_location_updates_scooter_timestamp ON location_updates(scooter_id, timestamp);

-- Drop the old trip_id column
ALTER TABLE location_updates DROP COLUMN trip_id;

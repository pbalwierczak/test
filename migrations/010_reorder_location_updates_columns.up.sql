-- Reorder columns in location_updates table to move scooter_id to second position
-- This requires recreating the table with the desired column order

-- Create a new table with the correct column order
CREATE TABLE location_updates_new (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scooter_id UUID NOT NULL REFERENCES scooters(id) ON DELETE CASCADE,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Copy data from old table to new table
INSERT INTO location_updates_new (id, scooter_id, latitude, longitude, timestamp, created_at, deleted_at)
SELECT id, scooter_id, latitude, longitude, timestamp, created_at, deleted_at
FROM location_updates;

-- Drop the old table
DROP TABLE location_updates;

-- Rename the new table to the original name
ALTER TABLE location_updates_new RENAME TO location_updates;

-- Recreate indexes
CREATE INDEX idx_location_updates_scooter_id ON location_updates(scooter_id);
CREATE INDEX idx_location_updates_timestamp ON location_updates(timestamp);
CREATE INDEX idx_location_updates_location ON location_updates(latitude, longitude);
CREATE INDEX idx_location_updates_scooter_timestamp ON location_updates(scooter_id, timestamp);
CREATE INDEX idx_location_updates_deleted_at ON location_updates(deleted_at);

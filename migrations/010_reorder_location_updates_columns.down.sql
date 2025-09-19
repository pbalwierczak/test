-- Revert column reordering in location_updates table
-- This reverts scooter_id back to its original position (after id)

-- Create a new table with the original column order
CREATE TABLE location_updates_old (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    scooter_id UUID NOT NULL REFERENCES scooters(id) ON DELETE CASCADE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Copy data from current table to old table
INSERT INTO location_updates_old (id, latitude, longitude, timestamp, created_at, scooter_id, deleted_at)
SELECT id, latitude, longitude, timestamp, created_at, scooter_id, deleted_at
FROM location_updates;

-- Drop the current table
DROP TABLE location_updates;

-- Rename the old table back to the original name
ALTER TABLE location_updates_old RENAME TO location_updates;

-- Recreate indexes
CREATE INDEX idx_location_updates_scooter_id ON location_updates(scooter_id);
CREATE INDEX idx_location_updates_timestamp ON location_updates(timestamp);
CREATE INDEX idx_location_updates_location ON location_updates(latitude, longitude);
CREATE INDEX idx_location_updates_scooter_timestamp ON location_updates(scooter_id, timestamp);
CREATE INDEX idx_location_updates_deleted_at ON location_updates(deleted_at);

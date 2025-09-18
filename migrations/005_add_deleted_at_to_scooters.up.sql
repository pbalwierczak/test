-- Add deleted_at column to scooters table for soft deletes
ALTER TABLE scooters ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE NULL;

-- Create index for deleted_at column
CREATE INDEX idx_scooters_deleted_at ON scooters(deleted_at);

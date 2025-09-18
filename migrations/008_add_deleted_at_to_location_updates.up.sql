-- Add deleted_at column to location_updates table for soft deletes
ALTER TABLE location_updates ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE NULL;

-- Create index for deleted_at column
CREATE INDEX idx_location_updates_deleted_at ON location_updates(deleted_at);

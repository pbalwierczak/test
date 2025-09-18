-- Add deleted_at column to trips table for soft deletes
ALTER TABLE trips ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE NULL;

-- Create index for deleted_at column
CREATE INDEX idx_trips_deleted_at ON trips(deleted_at);

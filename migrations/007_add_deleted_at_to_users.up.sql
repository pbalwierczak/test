-- Add deleted_at column to users table for soft deletes
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE NULL;

-- Create index for deleted_at column
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

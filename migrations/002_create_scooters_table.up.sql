-- Create scooters table
CREATE TABLE scooters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status VARCHAR(20) NOT NULL DEFAULT 'available' CHECK (status IN ('available', 'occupied')),
    current_latitude DECIMAL(10, 8) NOT NULL,
    current_longitude DECIMAL(11, 8) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);

-- Create indexes for performance
CREATE INDEX idx_scooters_status ON scooters(status);
CREATE INDEX idx_scooters_last_seen ON scooters(last_seen);
CREATE INDEX idx_scooters_location ON scooters(current_latitude, current_longitude);
CREATE INDEX idx_scooters_deleted_at ON scooters(deleted_at);

-- Add trigger to automatically update updated_at
CREATE TRIGGER update_scooters_updated_at 
    BEFORE UPDATE ON scooters 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

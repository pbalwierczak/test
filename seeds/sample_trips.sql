-- Seed data for trips and location_updates tables
-- This creates sample trips with location history for testing

-- Clean existing data (in reverse dependency order)
TRUNCATE TABLE location_updates CASCADE;
TRUNCATE TABLE trips CASCADE;

-- Sample completed trips
INSERT INTO trips (id, scooter_id, user_id, start_time, end_time, start_latitude, start_longitude, end_latitude, end_longitude, status, created_at, updated_at) VALUES
-- Ottawa trip 1 (completed)
('850e8400-e29b-41d4-a716-446655440001', '650e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 
 NOW() - INTERVAL '2 hours', NOW() - INTERVAL '1 hour 30 minutes', 45.4215, -75.6972, 45.4289, -75.6920, 'completed', 
 NOW() - INTERVAL '2 hours', NOW() - INTERVAL '1 hour 30 minutes'),

-- Ottawa trip 2 (completed)
('850e8400-e29b-41d4-a716-446655440002', '650e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', 
 NOW() - INTERVAL '3 hours', NOW() - INTERVAL '2 hours 30 minutes', 45.4200, -75.6950, 45.4250, -75.6850, 'completed', 
 NOW() - INTERVAL '3 hours', NOW() - INTERVAL '2 hours 30 minutes'),

-- Montreal trip 1 (completed)
('850e8400-e29b-41d4-a716-446655440003', '750e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440003', 
 NOW() - INTERVAL '4 hours', NOW() - INTERVAL '3 hours 30 minutes', 45.5017, -73.5673, 45.5089, -73.5620, 'completed', 
 NOW() - INTERVAL '4 hours', NOW() - INTERVAL '3 hours 30 minutes'),

-- Montreal trip 2 (completed)
('850e8400-e29b-41d4-a716-446655440004', '750e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440004', 
 NOW() - INTERVAL '5 hours', NOW() - INTERVAL '4 hours 30 minutes', 45.5000, -73.5650, 45.5150, -73.5550, 'completed', 
 NOW() - INTERVAL '5 hours', NOW() - INTERVAL '4 hours 30 minutes');

-- Location updates for completed trips
-- Trip 1 location updates (Ottawa)
INSERT INTO location_updates (id, trip_id, latitude, longitude, timestamp, created_at) VALUES
('950e8400-e29b-41d4-a716-446655440001', '850e8400-e29b-41d4-a716-446655440001', 45.4215, -75.6972, NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours'),
('950e8400-e29b-41d4-a716-446655440002', '850e8400-e29b-41d4-a716-446655440001', 45.4225, -75.6965, NOW() - INTERVAL '1 hour 57 minutes', NOW() - INTERVAL '1 hour 57 minutes'),
('950e8400-e29b-41d4-a716-446655440003', '850e8400-e29b-41d4-a716-446655440001', 45.4240, -75.6955, NOW() - INTERVAL '1 hour 54 minutes', NOW() - INTERVAL '1 hour 54 minutes'),
('950e8400-e29b-41d4-a716-446655440004', '850e8400-e29b-41d4-a716-446655440001', 45.4260, -75.6940, NOW() - INTERVAL '1 hour 51 minutes', NOW() - INTERVAL '1 hour 51 minutes'),
('950e8400-e29b-41d4-a716-446655440005', '850e8400-e29b-41d4-a716-446655440001', 45.4289, -75.6920, NOW() - INTERVAL '1 hour 30 minutes', NOW() - INTERVAL '1 hour 30 minutes');

-- Trip 2 location updates (Ottawa)
INSERT INTO location_updates (id, trip_id, latitude, longitude, timestamp, created_at) VALUES
('950e8400-e29b-41d4-a716-446655440006', '850e8400-e29b-41d4-a716-446655440002', 45.4200, -75.6950, NOW() - INTERVAL '3 hours', NOW() - INTERVAL '3 hours'),
('950e8400-e29b-41d4-a716-446655440007', '850e8400-e29b-41d4-a716-446655440002', 45.4210, -75.6940, NOW() - INTERVAL '2 hours 57 minutes', NOW() - INTERVAL '2 hours 57 minutes'),
('950e8400-e29b-41d4-a716-446655440008', '850e8400-e29b-41d4-a716-446655440002', 45.4225, -75.6925, NOW() - INTERVAL '2 hours 54 minutes', NOW() - INTERVAL '2 hours 54 minutes'),
('950e8400-e29b-41d4-a716-446655440009', '850e8400-e29b-41d4-a716-446655440002', 45.4240, -75.6900, NOW() - INTERVAL '2 hours 51 minutes', NOW() - INTERVAL '2 hours 51 minutes'),
('950e8400-e29b-41d4-a716-446655440010', '850e8400-e29b-41d4-a716-446655440002', 45.4250, -75.6850, NOW() - INTERVAL '2 hours 30 minutes', NOW() - INTERVAL '2 hours 30 minutes');

-- Trip 3 location updates (Montreal)
INSERT INTO location_updates (id, trip_id, latitude, longitude, timestamp, created_at) VALUES
('950e8400-e29b-41d4-a716-446655440011', '850e8400-e29b-41d4-a716-446655440003', 45.5017, -73.5673, NOW() - INTERVAL '4 hours', NOW() - INTERVAL '4 hours'),
('950e8400-e29b-41d4-a716-446655440012', '850e8400-e29b-41d4-a716-446655440003', 45.5030, -73.5660, NOW() - INTERVAL '3 hours 57 minutes', NOW() - INTERVAL '3 hours 57 minutes'),
('950e8400-e29b-41d4-a716-446655440013', '850e8400-e29b-41d4-a716-446655440003', 45.5050, -73.5645, NOW() - INTERVAL '3 hours 54 minutes', NOW() - INTERVAL '3 hours 54 minutes'),
('950e8400-e29b-41d4-a716-446655440014', '850e8400-e29b-41d4-a716-446655440003', 45.5070, -73.5630, NOW() - INTERVAL '3 hours 51 minutes', NOW() - INTERVAL '3 hours 51 minutes'),
('950e8400-e29b-41d4-a716-446655440015', '850e8400-e29b-41d4-a716-446655440003', 45.5089, -73.5620, NOW() - INTERVAL '3 hours 30 minutes', NOW() - INTERVAL '3 hours 30 minutes');

-- Trip 4 location updates (Montreal)
INSERT INTO location_updates (id, trip_id, latitude, longitude, timestamp, created_at) VALUES
('950e8400-e29b-41d4-a716-446655440016', '850e8400-e29b-41d4-a716-446655440004', 45.5000, -73.5650, NOW() - INTERVAL '5 hours', NOW() - INTERVAL '5 hours'),
('950e8400-e29b-41d4-a716-446655440017', '850e8400-e29b-41d4-a716-446655440004', 45.5020, -73.5630, NOW() - INTERVAL '4 hours 57 minutes', NOW() - INTERVAL '4 hours 57 minutes'),
('950e8400-e29b-41d4-a716-446655440018', '850e8400-e29b-41d4-a716-446655440004', 45.5050, -73.5600, NOW() - INTERVAL '4 hours 54 minutes', NOW() - INTERVAL '4 hours 54 minutes'),
('950e8400-e29b-41d4-a716-446655440019', '850e8400-e29b-41d4-a716-446655440004', 45.5080, -73.5570, NOW() - INTERVAL '4 hours 51 minutes', NOW() - INTERVAL '4 hours 51 minutes'),
('950e8400-e29b-41d4-a716-446655440020', '850e8400-e29b-41d4-a716-446655440004', 45.5150, -73.5550, NOW() - INTERVAL '4 hours 30 minutes', NOW() - INTERVAL '4 hours 30 minutes');


-- Seed data for trips and location_updates tables
-- This creates sample trips with location history for testing

-- Clean existing data (in reverse dependency order)
TRUNCATE TABLE location_updates CASCADE;
TRUNCATE TABLE trips CASCADE;

-- Sample completed trips
INSERT INTO trips (id, scooter_id, user_id, start_time, end_time, start_latitude, start_longitude, end_latitude, end_longitude, status, created_at, updated_at) VALUES
-- Ottawa trip 1 (completed)
('c3d4e5f6-a7b8-9012-cdef-345678901234', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', '550e8400-e29b-41d4-a716-446655440001', 
 NOW() - INTERVAL '2 hours', NOW() - INTERVAL '1 hour 30 minutes', 45.4215, -75.6972, 45.4289, -75.6920, 'completed', 
 NOW() - INTERVAL '2 hours', NOW() - INTERVAL '1 hour 30 minutes'),

-- Ottawa trip 2 (completed)
('c3d4e5f6-a7b8-9012-cdef-345678901235', 'a1b2c3d4-e5f6-7890-abcd-ef1234567891', '550e8400-e29b-41d4-a716-446655440002', 
 NOW() - INTERVAL '3 hours', NOW() - INTERVAL '2 hours 30 minutes', 45.4200, -75.6950, 45.4250, -75.6850, 'completed', 
 NOW() - INTERVAL '3 hours', NOW() - INTERVAL '2 hours 30 minutes'),

-- Montreal trip 1 (completed)
('c3d4e5f6-a7b8-9012-cdef-345678901236', 'b2c3d4e5-f6a7-8901-bcde-f23456789012', '550e8400-e29b-41d4-a716-446655440003', 
 NOW() - INTERVAL '4 hours', NOW() - INTERVAL '3 hours 30 minutes', 45.5017, -73.5673, 45.5089, -73.5620, 'completed', 
 NOW() - INTERVAL '4 hours', NOW() - INTERVAL '3 hours 30 minutes'),

-- Montreal trip 2 (completed)
('c3d4e5f6-a7b8-9012-cdef-345678901237', 'b2c3d4e5-f6a7-8901-bcde-f23456789013', '550e8400-e29b-41d4-a716-446655440004', 
 NOW() - INTERVAL '5 hours', NOW() - INTERVAL '4 hours 30 minutes', 45.5000, -73.5650, 45.5150, -73.5550, 'completed', 
 NOW() - INTERVAL '5 hours', NOW() - INTERVAL '4 hours 30 minutes');

-- Location updates for completed trips
-- Trip 1 location updates (Ottawa) - using scooter_id from trip 1
INSERT INTO location_updates (id, scooter_id, latitude, longitude, timestamp, created_at) VALUES
('d4e5f6a7-b8c9-0123-defa-456789012345', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 45.4215, -75.6972, NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours'),
('d4e5f6a7-b8c9-0123-defa-456789012346', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 45.4225, -75.6965, NOW() - INTERVAL '1 hour 57 minutes', NOW() - INTERVAL '1 hour 57 minutes'),
('d4e5f6a7-b8c9-0123-defa-456789012347', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 45.4240, -75.6955, NOW() - INTERVAL '1 hour 54 minutes', NOW() - INTERVAL '1 hour 54 minutes'),
('d4e5f6a7-b8c9-0123-defa-456789012348', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 45.4260, -75.6940, NOW() - INTERVAL '1 hour 51 minutes', NOW() - INTERVAL '1 hour 51 minutes'),
('d4e5f6a7-b8c9-0123-defa-456789012349', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 45.4215, -75.6972, NOW() - INTERVAL '1 hour 30 minutes', NOW() - INTERVAL '1 hour 30 minutes');

-- Trip 2 location updates (Ottawa) - using scooter_id from trip 2
INSERT INTO location_updates (id, scooter_id, latitude, longitude, timestamp, created_at) VALUES
('d4e5f6a7-b8c9-0123-defa-456789012350', 'a1b2c3d4-e5f6-7890-abcd-ef1234567891', 45.4200, -75.6950, NOW() - INTERVAL '3 hours', NOW() - INTERVAL '3 hours'),
('d4e5f6a7-b8c9-0123-defa-456789012351', 'a1b2c3d4-e5f6-7890-abcd-ef1234567891', 45.4210, -75.6940, NOW() - INTERVAL '2 hours 57 minutes', NOW() - INTERVAL '2 hours 57 minutes'),
('d4e5f6a7-b8c9-0123-defa-456789012352', 'a1b2c3d4-e5f6-7890-abcd-ef1234567891', 45.4225, -75.6925, NOW() - INTERVAL '2 hours 54 minutes', NOW() - INTERVAL '2 hours 54 minutes'),
('d4e5f6a7-b8c9-0123-defa-456789012353', 'a1b2c3d4-e5f6-7890-abcd-ef1234567891', 45.4240, -75.6900, NOW() - INTERVAL '2 hours 51 minutes', NOW() - INTERVAL '2 hours 51 minutes'),
('d4e5f6a7-b8c9-0123-defa-456789012354', 'a1b2c3d4-e5f6-7890-abcd-ef1234567891', 45.4200, -75.6950, NOW() - INTERVAL '2 hours 30 minutes', NOW() - INTERVAL '2 hours 30 minutes');

-- Trip 3 location updates (Montreal) - using scooter_id from trip 3
INSERT INTO location_updates (id, scooter_id, latitude, longitude, timestamp, created_at) VALUES
('d4e5f6a7-b8c9-0123-defa-456789012355', 'b2c3d4e5-f6a7-8901-bcde-f23456789012', 45.5017, -73.5673, NOW() - INTERVAL '4 hours', NOW() - INTERVAL '4 hours'),
('d4e5f6a7-b8c9-0123-defa-456789012356', 'b2c3d4e5-f6a7-8901-bcde-f23456789012', 45.5030, -73.5660, NOW() - INTERVAL '3 hours 57 minutes', NOW() - INTERVAL '3 hours 57 minutes'),
('d4e5f6a7-b8c9-0123-defa-456789012357', 'b2c3d4e5-f6a7-8901-bcde-f23456789012', 45.5050, -73.5645, NOW() - INTERVAL '3 hours 54 minutes', NOW() - INTERVAL '3 hours 54 minutes'),
('d4e5f6a7-b8c9-0123-defa-456789012358', 'b2c3d4e5-f6a7-8901-bcde-f23456789012', 45.5070, -73.5630, NOW() - INTERVAL '3 hours 51 minutes', NOW() - INTERVAL '3 hours 51 minutes'),
('d4e5f6a7-b8c9-0123-defa-456789012359', 'b2c3d4e5-f6a7-8901-bcde-f23456789012', 45.5017, -73.5673, NOW() - INTERVAL '3 hours 30 minutes', NOW() - INTERVAL '3 hours 30 minutes');

-- Trip 4 location updates (Montreal) - using scooter_id from trip 4
INSERT INTO location_updates (id, scooter_id, latitude, longitude, timestamp, created_at) VALUES
('d4e5f6a7-b8c9-0123-defa-456789012360', 'b2c3d4e5-f6a7-8901-bcde-f23456789013', 45.5000, -73.5650, NOW() - INTERVAL '5 hours', NOW() - INTERVAL '5 hours'),
('d4e5f6a7-b8c9-0123-defa-456789012361', 'b2c3d4e5-f6a7-8901-bcde-f23456789013', 45.5020, -73.5630, NOW() - INTERVAL '4 hours 57 minutes', NOW() - INTERVAL '4 hours 57 minutes'),
('d4e5f6a7-b8c9-0123-defa-456789012362', 'b2c3d4e5-f6a7-8901-bcde-f23456789013', 45.5050, -73.5600, NOW() - INTERVAL '4 hours 54 minutes', NOW() - INTERVAL '4 hours 54 minutes'),
('d4e5f6a7-b8c9-0123-defa-456789012363', 'b2c3d4e5-f6a7-8901-bcde-f23456789013', 45.5080, -73.5570, NOW() - INTERVAL '4 hours 51 minutes', NOW() - INTERVAL '4 hours 51 minutes'),
('d4e5f6a7-b8c9-0123-defa-456789012364', 'b2c3d4e5-f6a7-8901-bcde-f23456789013', 45.5000, -73.5650, NOW() - INTERVAL '4 hours 30 minutes', NOW() - INTERVAL '4 hours 30 minutes');


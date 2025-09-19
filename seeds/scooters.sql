-- Seed data for scooters table
-- This creates scooters distributed across Ottawa and Montreal areas

-- Clean existing data
TRUNCATE TABLE scooters CASCADE;

-- Ottawa area scooters (10 scooters)
INSERT INTO scooters (id, status, current_latitude, current_longitude, created_at, updated_at, last_seen) VALUES
-- Parliament Hill area
('650e8400-e29b-41d4-a716-446655440001', 'available', 45.4215, -75.6972, NOW() - INTERVAL '30 days', NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour'),
('650e8400-e29b-41d4-a716-446655440002', 'available', 45.4200, -75.6950, NOW() - INTERVAL '25 days', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours'),
('650e8400-e29b-41d4-a716-446655440003', 'available', 45.4230, -75.6990, NOW() - INTERVAL '20 days', NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '30 minutes'),

-- ByWard Market area
('650e8400-e29b-41d4-a716-446655440004', 'available', 45.4289, -75.6920, NOW() - INTERVAL '15 days', NOW() - INTERVAL '3 hours', NOW() - INTERVAL '3 hours'),
('650e8400-e29b-41d4-a716-446655440005', 'available', 45.4300, -75.6900, NOW() - INTERVAL '10 days', NOW() - INTERVAL '4 hours', NOW() - INTERVAL '4 hours'),
('650e8400-e29b-41d4-a716-446655440006', 'available', 45.4270, -75.6940, NOW() - INTERVAL '5 days', NOW() - INTERVAL '15 minutes', NOW() - INTERVAL '15 minutes'),

-- Rideau Centre area
('650e8400-e29b-41d4-a716-446655440007', 'available', 45.4250, -75.6850, NOW() - INTERVAL '3 days', NOW() - INTERVAL '5 hours', NOW() - INTERVAL '5 hours'),
('650e8400-e29b-41d4-a716-446655440008', 'available', 45.4260, -75.6830, NOW() - INTERVAL '2 days', NOW() - INTERVAL '6 hours', NOW() - INTERVAL '6 hours'),
('650e8400-e29b-41d4-a716-446655440009', 'available', 45.4240, -75.6870, NOW() - INTERVAL '1 day', NOW() - INTERVAL '45 minutes', NOW() - INTERVAL '45 minutes'),
('650e8400-e29b-41d4-a716-446655440010', 'available', 45.4280, -75.6810, NOW() - INTERVAL '12 hours', NOW() - INTERVAL '7 hours', NOW() - INTERVAL '7 hours'),

-- Montreal area scooters (10 scooters)
-- Old Port area
('750e8400-e29b-41d4-a716-446655440001', 'available', 45.5017, -73.5673, NOW() - INTERVAL '30 days', NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour'),
('750e8400-e29b-41d4-a716-446655440002', 'available', 45.5000, -73.5650, NOW() - INTERVAL '25 days', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours'),
('750e8400-e29b-41d4-a716-446655440003', 'available', 45.5030, -73.5690, NOW() - INTERVAL '20 days', NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '30 minutes'),

-- Downtown area
('750e8400-e29b-41d4-a716-446655440004', 'available', 45.5089, -73.5620, NOW() - INTERVAL '15 days', NOW() - INTERVAL '3 hours', NOW() - INTERVAL '3 hours'),
('750e8400-e29b-41d4-a716-446655440005', 'available', 45.5100, -73.5600, NOW() - INTERVAL '10 days', NOW() - INTERVAL '4 hours', NOW() - INTERVAL '4 hours'),
('750e8400-e29b-41d4-a716-446655440006', 'available', 45.5070, -73.5640, NOW() - INTERVAL '5 days', NOW() - INTERVAL '15 minutes', NOW() - INTERVAL '15 minutes'),

-- Plateau area
('750e8400-e29b-41d4-a716-446655440007', 'available', 45.5150, -73.5550, NOW() - INTERVAL '3 days', NOW() - INTERVAL '5 hours', NOW() - INTERVAL '5 hours'),
('750e8400-e29b-41d4-a716-446655440008', 'available', 45.5160, -73.5530, NOW() - INTERVAL '2 days', NOW() - INTERVAL '6 hours', NOW() - INTERVAL '6 hours'),
('750e8400-e29b-41d4-a716-446655440009', 'available', 45.5140, -73.5570, NOW() - INTERVAL '1 day', NOW() - INTERVAL '45 minutes', NOW() - INTERVAL '45 minutes'),
('750e8400-e29b-41d4-a716-446655440010', 'available', 45.5180, -73.5510, NOW() - INTERVAL '12 hours', NOW() - INTERVAL '7 hours', NOW() - INTERVAL '7 hours');

-- Seed data for scooters table
-- This creates scooters distributed across Ottawa and Montreal areas

-- Clean existing data
TRUNCATE TABLE scooters CASCADE;

-- Ottawa area scooters (10 scooters)
INSERT INTO scooters (id, status, current_latitude, current_longitude, created_at, updated_at, last_seen) VALUES
-- Parliament Hill area
('a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'available', 45.4215, -75.6972, NOW() - INTERVAL '30 days', NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour'),
('a1b2c3d4-e5f6-7890-abcd-ef1234567891', 'available', 45.4200, -75.6950, NOW() - INTERVAL '25 days', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours'),
('a1b2c3d4-e5f6-7890-abcd-ef1234567892', 'available', 45.4230, -75.6990, NOW() - INTERVAL '20 days', NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '30 minutes'),

-- ByWard Market area
('a1b2c3d4-e5f6-7890-abcd-ef1234567893', 'available', 45.4289, -75.6920, NOW() - INTERVAL '15 days', NOW() - INTERVAL '3 hours', NOW() - INTERVAL '3 hours'),
('a1b2c3d4-e5f6-7890-abcd-ef1234567894', 'available', 45.4300, -75.6900, NOW() - INTERVAL '10 days', NOW() - INTERVAL '4 hours', NOW() - INTERVAL '4 hours'),
('a1b2c3d4-e5f6-7890-abcd-ef1234567895', 'available', 45.4270, -75.6940, NOW() - INTERVAL '5 days', NOW() - INTERVAL '15 minutes', NOW() - INTERVAL '15 minutes'),

-- Rideau Centre area
('a1b2c3d4-e5f6-7890-abcd-ef1234567896', 'available', 45.4250, -75.6850, NOW() - INTERVAL '3 days', NOW() - INTERVAL '5 hours', NOW() - INTERVAL '5 hours'),
('a1b2c3d4-e5f6-7890-abcd-ef1234567897', 'available', 45.4260, -75.6830, NOW() - INTERVAL '2 days', NOW() - INTERVAL '6 hours', NOW() - INTERVAL '6 hours'),
('a1b2c3d4-e5f6-7890-abcd-ef1234567898', 'available', 45.4240, -75.6870, NOW() - INTERVAL '1 day', NOW() - INTERVAL '45 minutes', NOW() - INTERVAL '45 minutes'),
('a1b2c3d4-e5f6-7890-abcd-ef1234567899', 'available', 45.4280, -75.6810, NOW() - INTERVAL '12 hours', NOW() - INTERVAL '7 hours', NOW() - INTERVAL '7 hours'),

-- Montreal area scooters (10 scooters)
-- Old Port area
('b2c3d4e5-f6a7-8901-bcde-f23456789012', 'available', 45.5017, -73.5673, NOW() - INTERVAL '30 days', NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour'),
('b2c3d4e5-f6a7-8901-bcde-f23456789013', 'available', 45.5000, -73.5650, NOW() - INTERVAL '25 days', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours'),
('b2c3d4e5-f6a7-8901-bcde-f23456789014', 'available', 45.5030, -73.5690, NOW() - INTERVAL '20 days', NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '30 minutes'),

-- Downtown area
('b2c3d4e5-f6a7-8901-bcde-f23456789015', 'available', 45.5089, -73.5620, NOW() - INTERVAL '15 days', NOW() - INTERVAL '3 hours', NOW() - INTERVAL '3 hours'),
('b2c3d4e5-f6a7-8901-bcde-f23456789016', 'available', 45.5100, -73.5600, NOW() - INTERVAL '10 days', NOW() - INTERVAL '4 hours', NOW() - INTERVAL '4 hours'),
('b2c3d4e5-f6a7-8901-bcde-f23456789017', 'available', 45.5070, -73.5640, NOW() - INTERVAL '5 days', NOW() - INTERVAL '15 minutes', NOW() - INTERVAL '15 minutes'),

-- Plateau area
('b2c3d4e5-f6a7-8901-bcde-f23456789018', 'available', 45.5150, -73.5550, NOW() - INTERVAL '3 days', NOW() - INTERVAL '5 hours', NOW() - INTERVAL '5 hours'),
('b2c3d4e5-f6a7-8901-bcde-f23456789019', 'available', 45.5160, -73.5530, NOW() - INTERVAL '2 days', NOW() - INTERVAL '6 hours', NOW() - INTERVAL '6 hours'),
('b2c3d4e5-f6a7-8901-bcde-f23456789020', 'available', 45.5140, -73.5570, NOW() - INTERVAL '1 day', NOW() - INTERVAL '45 minutes', NOW() - INTERVAL '45 minutes'),
('b2c3d4e5-f6a7-8901-bcde-f23456789021', 'available', 45.5180, -73.5510, NOW() - INTERVAL '12 hours', NOW() - INTERVAL '7 hours', NOW() - INTERVAL '7 hours');

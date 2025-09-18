-- Seed data for users table
-- This creates test users for simulation and testing

INSERT INTO users (id, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440001', NOW() - INTERVAL '30 days', NOW() - INTERVAL '30 days'),
('550e8400-e29b-41d4-a716-446655440002', NOW() - INTERVAL '25 days', NOW() - INTERVAL '25 days'),
('550e8400-e29b-41d4-a716-446655440003', NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days'),
('550e8400-e29b-41d4-a716-446655440004', NOW() - INTERVAL '15 days', NOW() - INTERVAL '15 days'),
('550e8400-e29b-41d4-a716-446655440005', NOW() - INTERVAL '10 days', NOW() - INTERVAL '10 days'),
('550e8400-e29b-41d4-a716-446655440006', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days'),
('550e8400-e29b-41d4-a716-446655440007', NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days'),
('550e8400-e29b-41d4-a716-446655440008', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),
('550e8400-e29b-41d4-a716-446655440009', NOW() - INTERVAL '12 hours', NOW() - INTERVAL '12 hours'),
('550e8400-e29b-41d4-a716-446655440010', NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour');

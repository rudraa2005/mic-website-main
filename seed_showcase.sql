-- Use generate_uuid() or just let the database handle it if possible, but let's use explicit IDs for cross-referencing easily in script
-- We can use specific UUIDs for testing to avoid subquery failures if they don't exist

DELETE FROM work WHERE submission_id IN ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000003');
DELETE FROM submissions WHERE submission_id IN ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000003');
DELETE FROM users WHERE email = 'showcase@example.com';

-- Insert a dummy user
INSERT INTO users (id, name, email, password_hash, role) 
VALUES ('00000000-0000-0000-0000-000000000000', 'Showcase Student', 'showcase@example.com', 'password', 'STUDENT');

-- Insert dummy submissions
INSERT INTO submissions (submission_id, user_id, title, description, status, file_path)
VALUES 
    ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000000', 'Neural Network Optimizer', 'AI-based performance optimization for edge devices.', 'approved', 'https://images.unsplash.com/photo-1518770660439-4636190af475?auto=format&fit=crop&w=800&q=80'),
    ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000000', 'HydroHarvest', 'Atmospheric water generation for arid regions.', 'approved', 'https://images.unsplash.com/photo-1542332213-9b5a5a3fad35?auto=format&fit=crop&w=800&q=80'),
    ('00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000000', 'Sentinel Cyber', 'Decentralized intrusion detection system.', 'approved', 'https://images.unsplash.com/photo-1550751827-4bd374c3f58b?auto=format&fit=crop&w=800&q=80');

-- Insert dummy companies if they don't exist
INSERT INTO companies (id, name, logo_url)
VALUES 
    ('11111111-1111-1111-1111-111111111111', 'VentureX Capital', 'https://upload.wikimedia.org/wikipedia/commons/2/2f/Google_2015_logo.svg'),
    ('22222222-2222-2222-2222-222222222222', 'Future Labs', 'https://upload.wikimedia.org/wikipedia/commons/f/fa/Apple_logo_black.svg')
ON CONFLICT (id) DO NOTHING;

-- Map submissions to the work table
INSERT INTO work (submission_id, title, description, stage, company_id)
VALUES 
    ('00000000-0000-0000-0000-000000000001', 'Neural Network Optimizer', 'AI-based performance optimization for edge devices.', 'under_incubation', NULL),
    ('00000000-0000-0000-0000-000000000002', 'HydroHarvest', 'Atmospheric water generation for arid regions.', 'looking_for_funding', '11111111-1111-1111-1111-111111111111'),
    ('00000000-0000-0000-0000-000000000003', 'Sentinel Cyber', 'Decentralized intrusion detection system.', 'found_company', '22222222-2222-2222-2222-222222222222');

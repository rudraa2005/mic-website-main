-- Create admin user for MIC Website
-- Password should be changed after first login
-- Run this SQL in your PostgreSQL database

-- First, check if admin user already exists
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'admin@mic.mahe.edu.in') THEN
        INSERT INTO users (
            id,
            name,
            email,
            password_hash,
            role,
            created_at,
            updated_at
        ) VALUES (
            gen_random_uuid(),
            'MIC Admin',
            'admin@mic.mahe.edu.in',
            -- Password: MICAdmin@2024 (you should change this)
            -- This is a bcrypt hash - you may need to generate a new one
            '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
            'ADMIN',
            NOW(),
            NOW()
        );
        RAISE NOTICE 'Admin user created successfully';
    ELSE
        RAISE NOTICE 'Admin user already exists';
    END IF;
END $$;

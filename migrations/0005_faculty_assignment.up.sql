-- Migration: Add faculty assignment, tags, and work/incubation pipeline tables

-- Companies table for incubation partners
CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    logo_url TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Work items table for incubation pipeline
CREATE TABLE IF NOT EXISTS work (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID NOT NULL REFERENCES submissions(submission_id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    stage VARCHAR(50) DEFAULT 'under_incubation',
    progress_percent INT DEFAULT 0,
    company_id UUID REFERENCES companies(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Table for assigning faculty to submissions
CREATE TABLE IF NOT EXISTS submission_faculty (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID NOT NULL REFERENCES submissions(submission_id) ON DELETE CASCADE,
    faculty_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT NOW(),
    assigned_by UUID REFERENCES users(id),
    UNIQUE(submission_id, faculty_id)
);

-- Add tags/domain column to submissions if not exists
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'submissions' AND column_name = 'tags') THEN
        ALTER TABLE submissions ADD COLUMN tags TEXT[] DEFAULT '{}';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'submissions' AND column_name = 'domain') THEN
        ALTER TABLE submissions ADD COLUMN domain VARCHAR(255);
    END IF;
END $$;

-- Create indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_submission_faculty_faculty_id ON submission_faculty(faculty_id);
CREATE INDEX IF NOT EXISTS idx_submission_faculty_submission_id ON submission_faculty(submission_id);
CREATE INDEX IF NOT EXISTS idx_work_stage ON work(stage);
CREATE INDEX IF NOT EXISTS idx_work_submission_id ON work(submission_id);

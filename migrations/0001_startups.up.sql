-- Create table
CREATE TABLE startup_ideas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    stage TEXT NOT NULL DEFAULT 'draft',
    department TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Helpful indexes
CREATE INDEX idx_startup_owner ON startup_ideas(owner_id);
CREATE INDEX idx_startup_stage ON startup_ideas(stage);

-- Function to auto-update "updated_at"
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to run the function before every UPDATE
CREATE TRIGGER update_startup_ideas_updated_at
BEFORE UPDATE ON startup_ideas
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

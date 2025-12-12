CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    startup_id UUID NOT NULL,

    reviewer_id UUID NOT NULL,
    reviewer_name TEXT NOT NULL,
    reviewer_designation TEXT,

    rating NUMERIC(2,1),
    decision TEXT NOT NULL CHECK (decision IN ('approved', 'changes_requested', 'rejected')),
    summary TEXT,

    strengths JSONB DEFAULT '[]',
    recommendations JSONB DEFAULT '[]',

    created_at TIMESTAMP NOT NULL DEFAULT now()
);

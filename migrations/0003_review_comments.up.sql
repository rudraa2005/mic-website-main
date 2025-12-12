CREATE TABLE review_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    review_id UUID NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,

    commenter_id UUID NOT NULL,
    commenter_name TEXT NOT NULL,
    commenter_role TEXT, -- ex: "Professor", "Admin"

    comment TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

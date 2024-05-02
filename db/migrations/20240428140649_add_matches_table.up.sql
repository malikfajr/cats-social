CREATE TYPE STATUS_MATCH AS ENUM ('pending', 'approved', 'reject');

CREATE TABLE IF NOT EXISTS matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    issued_by JSONB NOT NULL,
    match_cat_detail JSONB NOT NULL,
    match_user_email VARCHAR(50) NOT NULL,
    user_cat_detail JSONB NOT NULL,
    message VARCHAR(150) NOT NULL,
    status STATUS_MATCH DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()
);
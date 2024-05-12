CREATE TYPE STATUS_MATCH AS ENUM ('pending', 'approved', 'reject');

CREATE TABLE IF NOT EXISTS matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    issued_user_id BIGINT NOT NULL,
    match_user_id BIGINT NOT NULL,
    user_cat_id UUID NOT NULL,
    match_cat_id UUID NOT NULL,
    message VARCHAR(150) NOT NULL,
    status STATUS_MATCH DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (issued_user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (match_user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (user_cat_id) REFERENCES cats(id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (match_cat_id) REFERENCES cats(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_match_created_at ON matches(created_at);

CREATE INDEX IF NOT EXISTS idx_match_status ON matches(status);
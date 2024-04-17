CREATE TABLE IF NOT EXISTS "matches" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    keyword_id UUID NOT NULL,
    source_url TEXT NOT NULL,
    content TEXT NOT NULL,
    found_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (keyword_id, source_url),
    FOREIGN KEY (keyword_id) REFERENCES keywords(id) ON DELETE CASCADE
);
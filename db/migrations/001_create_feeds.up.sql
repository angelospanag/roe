-- Create feeds table
CREATE TABLE IF NOT EXISTS feeds (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    url VARCHAR(512) NOT NULL UNIQUE,
    description TEXT,
    link VARCHAR(512),
    last_fetched_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_feeds_url ON feeds(url);
CREATE INDEX idx_feeds_last_fetched ON feeds(last_fetched_at);

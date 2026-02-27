-- Create posts table
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    feed_id INTEGER NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    title VARCHAR(512) NOT NULL,
    description TEXT,
    content TEXT,
    link VARCHAR(512) NOT NULL,
    author VARCHAR(255),
    published_at TIMESTAMP,
    guid VARCHAR(512) NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(feed_id, guid)
);

CREATE INDEX idx_posts_feed_id ON posts(feed_id);
CREATE INDEX idx_posts_is_read ON posts(is_read);
CREATE INDEX idx_posts_published_at ON posts(published_at);
CREATE INDEX idx_posts_guid ON posts(guid);

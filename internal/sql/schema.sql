CREATE TABLE feeds (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    name TEXT NOT NULL
);

CREATE TABLE feed_content (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    link TEXT NOT NULL,
    content TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    feed_id INTEGER NOT NULL,
    FOREIGN KEY (feed_id) REFERENCES feeds(id)
);

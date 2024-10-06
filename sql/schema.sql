CREATE TABLE feeds(
    id INTEGER PRIMARY KEY,
    url TEXT NOT NULL,
    name TEXT NOT NULL
);
CREATE TABLE feed_content(
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    link TEXT NOT NULL,
    content TEXT NOT NULL,
    is_read INTEGER NOT NULL DEFAULT 0,
    feed_id INTEGER NOT NULL,
    FOREIGN KEY (feed_id) REFERENCES feeds(id)
);
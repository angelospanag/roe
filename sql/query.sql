-- name: GetFeeds :many
SELECT feeds.id,
    feeds.url,
    feeds.name,
    COUNT(feed_content.id) as unread_items_count
FROM feeds
    LEFT JOIN feed_content ON feeds.id = feed_content.feed_id
WHERE feed_content.is_read = 0
ORDER BY name;
-- name: GetFeedItems :many
SELECT id,
    title,
    description,
    link,
    content,
    is_read
FROM feed_content
WHERE feed_id = ?;
-- name: GetFeedItem :one
SELECT id,
    title,
    description,
    link,
    content,
    is_read
FROM feed_content
WHERE feed_id = ?
    and id = ?;
-- name: AddFeed :one
INSERT INTO feeds (url, name)
VALUES (?, ?)
RETURNING *;
-- name: AddFeedContent :exec
INSERT INTO feed_content (title, description, link, content, feed_id)
VALUES(?, ?, ?, ?, ?);
-- name: UpdateFeedItem :one
UPDATE feed_content
SET is_read = ?
WHERE feed_id = ?
    AND id = ?
RETURNING *;
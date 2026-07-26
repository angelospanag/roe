-- name: CreatePost :one
INSERT INTO posts (feed_id, title, description, content, link, author, published_at, guid)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (feed_id, guid) DO UPDATE
SET title = EXCLUDED.title,
    description = EXCLUDED.description,
    content = EXCLUDED.content,
    link = EXCLUDED.link,
    author = EXCLUDED.author,
    published_at = EXCLUDED.published_at,
    updated_at = NOW()
RETURNING *;

-- name: GetPost :one
SELECT * FROM posts
WHERE id = $1;

-- name: ListPosts :many
SELECT * FROM posts
ORDER BY published_at DESC NULLS LAST
LIMIT $1 OFFSET $2;

-- name: ListPostsByFeed :many
SELECT * FROM posts
WHERE feed_id = $1
ORDER BY published_at DESC NULLS LAST
LIMIT $2 OFFSET $3;

-- name: ListUnreadPosts :many
SELECT * FROM posts
WHERE is_read = FALSE
ORDER BY published_at DESC NULLS LAST
LIMIT $1 OFFSET $2;

-- name: ListUnreadPostsByFeed :many
SELECT * FROM posts
WHERE feed_id = $1 AND is_read = FALSE
ORDER BY published_at DESC NULLS LAST
LIMIT $2 OFFSET $3;

-- name: MarkPostAsRead :exec
UPDATE posts
SET is_read = TRUE, updated_at = NOW()
WHERE id = $1;

-- name: MarkPostAsUnread :exec
UPDATE posts
SET is_read = FALSE, updated_at = NOW()
WHERE id = $1;

-- name: MarkAllPostsAsRead :exec
UPDATE posts
SET is_read = TRUE, updated_at = NOW()
WHERE feed_id = $1;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;

-- name: CountUnreadPosts :one
SELECT COUNT(*) FROM posts
WHERE is_read = FALSE;

-- name: CountUnreadPostsByFeed :one
SELECT COUNT(*) FROM posts
WHERE feed_id = $1 AND is_read = FALSE;

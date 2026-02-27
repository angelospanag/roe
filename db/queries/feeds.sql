-- name: CreateFeed :one
INSERT INTO feeds (title, url, description, link)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetFeed :one
SELECT * FROM feeds
WHERE id = $1;

-- name: GetFeedByURL :one
SELECT * FROM feeds
WHERE url = $1;

-- name: ListFeeds :many
SELECT * FROM feeds
ORDER BY created_at DESC;

-- name: UpdateFeed :one
UPDATE feeds
SET title = $2, description = $3, link = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateFeedLastFetched :exec
UPDATE feeds
SET last_fetched_at = NOW()
WHERE id = $1;

-- name: DeleteFeed :exec
DELETE FROM feeds
WHERE id = $1;

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package db

import (
	"context"
)

const addFeed = `-- name: AddFeed :one
INSERT INTO feeds (url, name)
VALUES ($1, $2)
RETURNING id, url, name
`

type AddFeedParams struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

func (q *Queries) AddFeed(ctx context.Context, arg AddFeedParams) (Feed, error) {
	row := q.db.QueryRow(ctx, addFeed, arg.Url, arg.Name)
	var i Feed
	err := row.Scan(&i.ID, &i.Url, &i.Name)
	return i, err
}

const addFeedContent = `-- name: AddFeedContent :exec
INSERT INTO feed_content (title, description, link, content, feed_id)
VALUES($1, $2, $3, $4, $5)
`

type AddFeedContentParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Content     string `json:"content"`
	FeedID      int32  `json:"feed_id"`
}

func (q *Queries) AddFeedContent(ctx context.Context, arg AddFeedContentParams) error {
	_, err := q.db.Exec(ctx, addFeedContent,
		arg.Title,
		arg.Description,
		arg.Link,
		arg.Content,
		arg.FeedID,
	)
	return err
}

const getFeedItem = `-- name: GetFeedItem :one
SELECT id,
    title,
    description,
    link,
    content,
    is_read
FROM feed_content
WHERE feed_id = $1
    and id = $2
`

type GetFeedItemParams struct {
	FeedID int32 `json:"feed_id"`
	ID     int32 `json:"id"`
}

type GetFeedItemRow struct {
	ID          int32  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Content     string `json:"content"`
	IsRead      bool   `json:"is_read"`
}

func (q *Queries) GetFeedItem(ctx context.Context, arg GetFeedItemParams) (GetFeedItemRow, error) {
	row := q.db.QueryRow(ctx, getFeedItem, arg.FeedID, arg.ID)
	var i GetFeedItemRow
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Link,
		&i.Content,
		&i.IsRead,
	)
	return i, err
}

const getFeedItems = `-- name: GetFeedItems :many
SELECT id,
    title,
    description,
    link,
    content,
    is_read
FROM feed_content
WHERE feed_id = $1
`

type GetFeedItemsRow struct {
	ID          int32  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Content     string `json:"content"`
	IsRead      bool   `json:"is_read"`
}

func (q *Queries) GetFeedItems(ctx context.Context, feedID int32) ([]GetFeedItemsRow, error) {
	rows, err := q.db.Query(ctx, getFeedItems, feedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedItemsRow
	for rows.Next() {
		var i GetFeedItemsRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Link,
			&i.Content,
			&i.IsRead,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFeeds = `-- name: GetFeeds :many
SELECT feeds.id,
    feeds.url,
    feeds.name,
    COUNT(feed_content.id) as unread_items_count
FROM feeds
    LEFT JOIN feed_content ON feeds.id = feed_content.feed_id
WHERE feed_content.is_read = FALSE
GROUP BY feeds.id
ORDER BY name
`

type GetFeedsRow struct {
	ID               int32  `json:"id"`
	Url              string `json:"url"`
	Name             string `json:"name"`
	UnreadItemsCount int64  `json:"unread_items_count"`
}

func (q *Queries) GetFeeds(ctx context.Context) ([]GetFeedsRow, error) {
	rows, err := q.db.Query(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedsRow
	for rows.Next() {
		var i GetFeedsRow
		if err := rows.Scan(
			&i.ID,
			&i.Url,
			&i.Name,
			&i.UnreadItemsCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateFeedItem = `-- name: UpdateFeedItem :one
UPDATE feed_content
SET is_read = $1
WHERE feed_id = $2
    AND id = $3
RETURNING id, title, description, link, content, is_read, feed_id
`

type UpdateFeedItemParams struct {
	IsRead bool  `json:"is_read"`
	FeedID int32 `json:"feed_id"`
	ID     int32 `json:"id"`
}

func (q *Queries) UpdateFeedItem(ctx context.Context, arg UpdateFeedItemParams) (FeedContent, error) {
	row := q.db.QueryRow(ctx, updateFeedItem, arg.IsRead, arg.FeedID, arg.ID)
	var i FeedContent
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Link,
		&i.Content,
		&i.IsRead,
		&i.FeedID,
	)
	return i, err
}

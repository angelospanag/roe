// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

type Feed struct {
	ID   int64  `json:"id"`
	Url  string `json:"url"`
	Name string `json:"name"`
}

type FeedContent struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Content     string `json:"content"`
	IsRead      int64  `json:"is_read"`
	FeedID      int64  `json:"feed_id"`
}

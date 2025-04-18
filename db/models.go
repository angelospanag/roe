// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

type Feed struct {
	ID   int32  `json:"id"`
	Url  string `json:"url"`
	Name string `json:"name"`
}

type FeedContent struct {
	ID          int32  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Content     string `json:"content"`
	IsRead      bool   `json:"is_read"`
	FeedID      int32  `json:"feed_id"`
}

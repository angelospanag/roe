package post

import "time"

// PostResponse represents a post in API responses
type PostResponse struct {
	ID          int32      `json:"id" example:"1"`
	FeedID      int32      `json:"feed_id" example:"1"`
	Title       string     `json:"title" example:"Breaking News"`
	Description *string    `json:"description,omitempty"`
	Content     *string    `json:"content,omitempty"`
	Link        string     `json:"link" example:"https://example.com/article"`
	Author      *string    `json:"author,omitempty" example:"John Doe"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	GUID        string     `json:"guid" example:"unique-guid-123"`
	IsRead      bool       `json:"is_read" example:"false"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// MarkReadRequest represents a request to mark posts as read/unread
type MarkReadRequest struct {
	IsRead bool `json:"is_read" example:"true"`
}

// ListParams represents pagination parameters
type ListParams struct {
	Limit  int32 `query:"limit" minimum:"1" maximum:"100" default:"20" example:"20"`
	Offset int32 `query:"offset" minimum:"0" default:"0" example:"0"`
}

// ListParamsWithFilters includes pagination and optional feed filter
type ListParamsWithFilters struct {
	ListParams
	FeedID     int32 `query:"feed_id" minimum:"0" default:"0" example:"1" doc:"Filter by feed ID (0 = all feeds)"`
	UnreadOnly bool  `query:"unread_only" default:"false" example:"false" doc:"Show only unread posts"`
}

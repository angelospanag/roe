package feed

import "time"

// Input represents the input for creating/updating a feed
type Input struct {
	Title       string `json:"title" maxLength:"255" example:"TechCrunch"`
	URL         string `json:"url" maxLength:"512" format:"uri" example:"https://techcrunch.com/feed/"`
	Description string `json:"description,omitempty" example:"Technology news and analysis"`
	Link        string `json:"link,omitempty" format:"uri" example:"https://techcrunch.com"`
}

// FeedResponse represents a feed in API responses
type FeedResponse struct {
	ID            int32      `json:"id" example:"1"`
	Title         string     `json:"title" example:"TechCrunch"`
	URL           string     `json:"url" example:"https://techcrunch.com/feed/"`
	Description   *string    `json:"description,omitempty"`
	Link          *string    `json:"link,omitempty"`
	LastFetchedAt *time.Time `json:"last_fetched_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// RefreshRequest represents a request to refresh feeds
type RefreshRequest struct {
	FeedID *int32 `json:"feed_id,omitempty" example:"1" doc:"Optional feed ID to refresh specific feed"`
}

// RefreshResponse represents the result of a refresh operation
type RefreshResponse struct {
	Message      string `json:"message" example:"Feeds refreshed successfully"`
	FeedsUpdated int    `json:"feeds_updated" example:"5"`
	PostsAdded   int    `json:"posts_added" example:"42"`
}

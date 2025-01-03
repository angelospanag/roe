package internal

import (
	"context"

	"github.com/angelospanag/rss-llm-go/db"
	"github.com/mmcdole/gofeed"
)

// GetAllFeeds stores a new feed in the database
func AddFeed(queries *db.Queries, url string, name *string) (*db.Feed, error) {

	// Parse feed, optionally add a custom feed name if one is passed
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)
	if name == nil {
		name = &feed.Title
	}

	// Add feed
	newFeed, err := queries.AddFeed(context.Background(), db.AddFeedParams{
		Url:  url,
		Name: *name,
	})
	if err != nil {
		return nil, err
	}

	// Add feed content
	for _, item := range feed.Items {
		err = queries.AddFeedContent(context.Background(), db.AddFeedContentParams{
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
			Content:     item.Content,
			FeedID:      newFeed.ID})
		if err != nil {
			return nil, err
		}
	}
	return &newFeed, nil
}

// GetAllFeeds retrieves the all the feeds from the database
func GetAllFeeds(queries *db.Queries) (*[]db.GetFeedsRow, error) {
	feeds, err := queries.GetFeeds(context.Background())
	if err != nil {
		return nil, err
	}
	return &feeds, nil
}

// GetAllFeeds retrieves the all the items of a feed from the database
func GetFeedItems(queries *db.Queries, feedID int32) (*[]db.GetFeedItemsRow, error) {

	feedItems, err := queries.GetFeedItems(context.Background(), feedID)
	if err != nil {
		return nil, err
	}
	return &feedItems, nil
}

package feeds

import (
	"context"

	"github.com/angelospanag/roe/internal/app"
	"github.com/angelospanag/roe/internal/db"
	"github.com/mmcdole/gofeed"
)

// Service is a struct that provides methods to interact with the database
type Service struct {
	app *app.App
}

// NewService creates a new Service instance
func NewService(app *app.App) *Service {
	return &Service{
		app: app,
	}
}

// AddFeed adds a new feed to the database and populates it with the feed content
func (s *Service) AddFeed(
	ctx context.Context,
	queries *db.Queries,
	url string,
	name *string,
) (*db.Feed, error) {
	// Parse feed, optionally add a custom feed name if one is passed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}

	if name == nil {
		name = &feed.Title
	}

	// Add feed
	newFeed, err := queries.AddFeed(ctx, db.AddFeedParams{
		Url:  url,
		Name: *name,
	})
	if err != nil {
		return nil, err
	}

	// Add feed content
	for _, item := range feed.Items {
		err = queries.AddFeedContent(ctx, db.AddFeedContentParams{
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
			Content:     item.Content,
			FeedID:      newFeed.ID,
		})
		if err != nil {
			return nil, err
		}
	}
	return &newFeed, nil
}

// GetAllFeeds retrieves all the feeds from the database
func (s *Service) GetAllFeeds(ctx context.Context, queries *db.Queries) ([]db.GetFeedsRow, error) {
	feeds, err := queries.GetFeeds(ctx)
	if err != nil {
		return nil, err
	}
	return feeds, nil
}

// GetFeedItems retrieves all feed items from the database
func (s *Service) GetFeedItems(
	ctx context.Context,
	queries *db.Queries,
	feedID int32,
) (*[]db.GetFeedItemsRow, error) {
	feedItems, err := queries.GetFeedItems(ctx, feedID)
	if err != nil {
		return nil, err
	}
	return &feedItems, nil
}

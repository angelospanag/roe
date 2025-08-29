package feeds

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/angelospanag/roe/internal/app"
	"github.com/angelospanag/roe/internal/db"
	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers the routes for the feeds feature.
func RegisterRoutes(api huma.API, app *app.App) {
	feedsService := NewService(app)

	type GetFeedsOutput struct {
		Body struct {
			Feeds []db.GetFeedsRow `json:"feeds"`
		}
	}
	// Register GET /feeds
	huma.Register(api, huma.Operation{
		OperationID:   "get-feeds",
		Method:        http.MethodGet,
		Path:          "/feeds",
		Summary:       "Get all feeds",
		Tags:          []string{"Feeds"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *struct{}) (*GetFeedsOutput, error) {
		resp := &GetFeedsOutput{}
		feeds, err := feedsService.GetAllFeeds(ctx, app.DBQueries)
		if err != nil {
			slog.Error("error fetching feeds", "error", err.Error())
			return nil, huma.Error500InternalServerError("something went wrong")
		}

		resp.Body.Feeds = feeds

		return resp, nil
	})

	// Create a feed
	type AddFeedInput struct {
		Body struct {
			Url  string `json:"url" type:"uri" example:"https://www.skai.gr/feed.xml"`
			Name string `json:"name,omitempty" maxLength:"80" example:"SKAI News"`
		}
	}

	type AddFeedOutput struct {
		Body struct {
			Feed db.Feed `json:"feed"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID:   "post-feed",
		Method:        http.MethodPost,
		Path:          "/feeds",
		Summary:       "Create a feed",
		Tags:          []string{"Feeds"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *AddFeedInput) (*AddFeedOutput, error) {
		resp := &AddFeedOutput{}

		var feedName *string
		if i.Body.Name == "" {
			feedName = nil
		} else {
			feedName = &i.Body.Name
		}

		newFeed, err := feedsService.AddFeed(ctx, app.DBQueries, i.Body.Url, feedName)
		if err != nil {
			slog.Error("error adding feeds", "error", err.Error())
			return nil, huma.Error500InternalServerError("Something went wrong")
		}
		resp.Body.Feed = *newFeed
		return resp, nil
	})

	// Get a feed's items
	type GetFeedItemsInput struct {
		FeedID int32 `path:"feedID"`
	}

	type GetFeedItemsOutput struct {
		Body struct {
			Items []db.GetFeedItemsRow `json:"items"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID:   "get-feed-items",
		Method:        http.MethodGet,
		Path:          "/feeds/{feedID}/items",
		Summary:       "Get feed items",
		Tags:          []string{"Feeds, Item"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *GetFeedItemsInput) (*GetFeedItemsOutput, error) {
		resp := &GetFeedItemsOutput{}
		feedItems, err := feedsService.GetFeedItems(ctx, app.DBQueries, i.FeedID)
		if err != nil {
			slog.Error("error getting feed items", "error", err.Error())
			return nil, huma.Error500InternalServerError("something went wrong")
		}

		resp.Body.Items = *feedItems

		return resp, nil
	})

	// Get a feed item
	type GetFeedItemInput struct {
		ItemID int32 `path:"itemID"`
		FeedID int32 `path:"feedID"`
	}

	type GetFeedItemOutput struct {
		Body struct {
			Item db.GetFeedItemRow `json:"item"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID:   "get-feed-item",
		Method:        http.MethodGet,
		Path:          "/feeds/{feedID}/items/{itemID}",
		Summary:       "Get feed item",
		Tags:          []string{"Feeds, Item"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *GetFeedItemInput) (*GetFeedItemOutput, error) {
		resp := &GetFeedItemOutput{}
		feedItem, err := app.DBQueries.GetFeedItem(ctx, db.GetFeedItemParams{
			FeedID: i.FeedID,
			ID:     i.ItemID,
		})
		if err != nil {
			slog.Error("error getting feed item", "error", err.Error())
			return nil, huma.Error500InternalServerError("something went wrong")
		}

		resp.Body.Item = feedItem

		return resp, nil
	})

	// Update feed item
	type UpdateFeedItemInput struct {
		FeedID int32 `path:"feedID"`
		ItemID int32 `path:"itemID"`
		Body   struct {
			IsRead bool `json:"is_read"`
		}
	}

	type UpdateFeedItemOutput struct {
		Body struct {
			Item db.FeedContent `json:"item"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID:   "update-feed-item",
		Method:        http.MethodPost,
		Path:          "/feeds/{feedID}/items/{itemID}",
		Summary:       "Update feed item",
		Tags:          []string{"Feeds, Item"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *UpdateFeedItemInput) (*UpdateFeedItemOutput, error) {
		resp := &UpdateFeedItemOutput{}

		updatedFeedItem, err := app.DBQueries.UpdateFeedItem(
			ctx,
			db.UpdateFeedItemParams{
				FeedID: i.FeedID,
				ID:     i.ItemID,
				IsRead: i.Body.IsRead,
			},
		)
		if err != nil {
			slog.Error("error updating feed item", "error", err.Error())
			return nil, huma.Error500InternalServerError("Something went wrong")
		}
		resp.Body.Item = updatedFeedItem
		return resp, nil
	})
}

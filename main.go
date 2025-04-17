package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/angelospanag/rss-llm-go/db"
	"github.com/angelospanag/rss-llm-go/internal"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
)

//go:embed frontend/dist/*
var spaFiles embed.FS

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "user=user host=localhost password=password dbname=testdb")
	if err != nil {
		log.Fatalf("Error opening database connection %v", err.Error())
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	// Create a new router & API
	router := chi.NewMux()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

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
		feeds, err := internal.GetAllFeeds(queries)
		if err != nil {
			slog.Error("error fetching feeds", "error", err.Error())
			return nil, huma.Error500InternalServerError("something went wrong")
		}

		resp.Body.Feeds = *feeds

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

		newFeed, err := internal.AddFeed(queries, i.Body.Url, feedName)
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
		feedItems, err := internal.GetFeedItems(queries, i.FeedID)
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
		feedItem, err := queries.GetFeedItem(context.Background(), db.GetFeedItemParams{
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

		updatedFeedItem, err := queries.UpdateFeedItem(context.Background(), db.UpdateFeedItemParams{
			FeedID: i.FeedID,
			ID:     i.ItemID,
			IsRead: i.Body.IsRead,
		})

		if err != nil {
			slog.Error("error updating feed item", "error", err.Error())
			return nil, huma.Error500InternalServerError("Something went wrong")
		}
		resp.Body.Item = updatedFeedItem
		return resp, nil
	})

	router.Handle("/*", SPAHandler())

	// Start server
	http.ListenAndServe("127.0.0.1:8000", router)
}

// https://github.com/go-chi/chi/issues/611
func SPAHandler() http.HandlerFunc {
	spaFS, err := fs.Sub(spaFiles, "frontend/dist")
	if err != nil {
		panic(fmt.Errorf("failed getting the sub tree for the site files: %w", err))
	}
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := spaFS.Open(strings.TrimPrefix(path.Clean(r.URL.Path), "/"))
		if err == nil {
			defer f.Close()
		}
		if os.IsNotExist(err) {
			r.URL.Path = "/"
		}
		http.FileServer(http.FS(spaFS)).ServeHTTP(w, r)
	}
}

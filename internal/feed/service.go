package feed

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/angelospanag/roe/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmcdole/gofeed"
)

// Service handles RSS feed operations
type Service struct {
	queries *db.Queries
	pool    *pgxpool.Pool
	parser  *gofeed.Parser
	logger  *slog.Logger
}

// NewService creates a new feed service
func NewService(pool *pgxpool.Pool, logger *slog.Logger) *Service {
	return &Service{
		queries: db.New(pool),
		pool:    pool,
		parser:  gofeed.NewParser(),
		logger:  logger,
	}
}

// RefreshFeed fetches and updates posts for a specific feed
func (s *Service) RefreshFeed(ctx context.Context, feedID int32) (int, error) {
	s.logger.Info("refreshing feed", "feed_id", feedID)

	// Get feed details
	feed, err := s.queries.GetFeed(ctx, feedID)
	if err != nil {
		return 0, fmt.Errorf("failed to get feed: %w", err)
	}

	// Parse RSS feed
	parsedFeed, err := s.parser.ParseURLWithContext(feed.Url, ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to parse feed: %w", err)
	}

	// Update feed metadata
	_, err = s.queries.UpdateFeed(ctx, db.UpdateFeedParams{
		ID: feedID,
		Title: func() string {
			if parsedFeed.Title != "" {
				return parsedFeed.Title
			}
			return feed.Title
		}(),
		Description: pgtype.Text(sql.NullString{
			String: parsedFeed.Description,
			Valid:  parsedFeed.Description != "",
		}),
		Link: pgtype.Text(sql.NullString{
			String: parsedFeed.Link,
			Valid:  parsedFeed.Link != "",
		}),
	})
	if err != nil {
		s.logger.Error("failed to update feed metadata", "error", err)
	}

	// Add/update posts
	postsAdded := 0
	for _, item := range parsedFeed.Items {
		guid := item.GUID
		if guid == "" {
			guid = item.Link
		}

		publishedAt := pgtype.Timestamp{Valid: false}
		if item.PublishedParsed != nil {
			publishedAt = pgtype.Timestamp{
				Time:  *item.PublishedParsed,
				Valid: true,
			}
		}

		_, err := s.queries.CreatePost(ctx, db.CreatePostParams{
			FeedID: feedID,
			Title:  item.Title,
			Description: pgtype.Text(sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			}),
			Content: pgtype.Text(sql.NullString{
				String: item.Content,
				Valid:  item.Content != "",
			}),
			Link: item.Link,
			Author: pgtype.Text(sql.NullString{
				String: func() string {
					if item.Author != nil {
						return item.Author.Name
					}
					return ""
				}(),
				Valid: item.Author != nil && item.Author.Name != "",
			}),
			PublishedAt: publishedAt,
			Guid:        guid,
		})

		if err != nil {
			s.logger.Warn("failed to create post", "error", err, "guid", guid)
		} else {
			postsAdded++
		}
	}

	// Update last fetched timestamp
	err = s.queries.UpdateFeedLastFetched(ctx, feedID)
	if err != nil {
		s.logger.Error("failed to update last fetched timestamp", "error", err)
	}

	s.logger.Info("feed refreshed", "feed_id", feedID, "posts_added", postsAdded)
	return postsAdded, nil
}

// RefreshAllFeeds fetches and updates posts for all feeds
func (s *Service) RefreshAllFeeds(ctx context.Context) (int, int, error) {
	s.logger.Info("refreshing all feeds")

	feeds, err := s.queries.ListFeeds(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to list feeds: %w", err)
	}

	feedsUpdated := 0
	totalPostsAdded := 0

	for _, feed := range feeds {
		postsAdded, err := s.RefreshFeed(ctx, feed.ID)
		if err != nil {
			s.logger.Error("failed to refresh feed", "feed_id", feed.ID, "error", err)
			continue
		}
		feedsUpdated++
		totalPostsAdded += postsAdded
	}

	s.logger.Info("all feeds refreshed", "feeds_updated", feedsUpdated, "posts_added", totalPostsAdded)
	return feedsUpdated, totalPostsAdded, nil
}

// GetQueries returns the database queries instance
func (s *Service) GetQueries() *db.Queries {
	return s.queries
}

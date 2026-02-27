package post

import (
	"context"
	"log/slog"

	"github.com/angelospanag/roe/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service handles post operations
type Service struct {
	queries *db.Queries
	logger  *slog.Logger
}

// NewService creates a new post service
func NewService(pool *pgxpool.Pool, logger *slog.Logger) *Service {
	return &Service{
		queries: db.New(pool),
		logger:  logger,
	}
}

// GetQueries returns the database queries interface
func (s *Service) GetQueries() *db.Queries {
	return s.queries
}

// MarkAllPostsAsRead marks all posts in a feed as read
func (s *Service) MarkAllPostsAsRead(ctx context.Context, feedID int32) error {
	return s.queries.MarkAllPostsAsRead(ctx, feedID)
}

// CountUnreadPostsByFeed returns the count of unread posts for a feed
func (s *Service) CountUnreadPostsByFeed(ctx context.Context, feedID int32) (int64, error) {
	return s.queries.CountUnreadPostsByFeed(ctx, feedID)
}

// ListPostsByFeed lists posts for a specific feed
func (s *Service) ListPostsByFeed(ctx context.Context, params db.ListPostsByFeedParams) ([]db.Post, error) {
	return s.queries.ListPostsByFeed(ctx, params)
}

// ListUnreadPostsByFeed lists unread posts for a specific feed
func (s *Service) ListUnreadPostsByFeed(ctx context.Context, params db.ListUnreadPostsByFeedParams) ([]db.Post, error) {
	return s.queries.ListUnreadPostsByFeed(ctx, params)
}

// ListPosts lists all posts
func (s *Service) ListPosts(ctx context.Context, params db.ListPostsParams) ([]db.Post, error) {
	return s.queries.ListPosts(ctx, params)
}

// ListUnreadPosts lists all unread posts
func (s *Service) ListUnreadPosts(ctx context.Context, params db.ListUnreadPostsParams) ([]db.Post, error) {
	return s.queries.ListUnreadPosts(ctx, params)
}

// GetPost retrieves a specific post by ID
func (s *Service) GetPost(ctx context.Context, id int32) (db.Post, error) {
	return s.queries.GetPost(ctx, id)
}

// MarkPostAsRead marks a post as read
func (s *Service) MarkPostAsRead(ctx context.Context, postID int32) error {
	return s.queries.MarkPostAsRead(ctx, postID)
}

// MarkPostAsUnread marks a post as unread
func (s *Service) MarkPostAsUnread(ctx context.Context, postID int32) error {
	return s.queries.MarkPostAsUnread(ctx, postID)
}

// CountUnreadPosts returns the total count of unread posts
func (s *Service) CountUnreadPosts(ctx context.Context) (int64, error) {
	return s.queries.CountUnreadPosts(ctx)
}

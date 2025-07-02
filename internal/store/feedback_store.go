package store

import (
	"context"
	"time"

	"github.com/y3933y3933/joker/internal/db/sqlc"
)

type Feedback struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type PostgresFeedbackStore struct {
	queries *sqlc.Queries
}

func NewPostgresFeedStore(queries *sqlc.Queries) *PostgresFeedbackStore {
	return &PostgresFeedbackStore{queries: queries}
}

type FeedbackStore interface {
	Create(context.Context, *Feedback) error
}

func (pg *PostgresFeedbackStore) Create(ctx context.Context, feedback *Feedback) error {
	return pg.queries.CreateFeedback(ctx, sqlc.CreateFeedbackParams{
		Type:    feedback.Type,
		Content: feedback.Content,
	})
}

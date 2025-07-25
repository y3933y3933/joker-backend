package store

import (
	"context"
	"time"

	"github.com/y3933y3933/joker/internal/db/sqlc"
)

type Feedback struct {
	ID           int64     `json:"id"`
	Type         string    `json:"type"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"createdAt"`
	ReviewStatus string    `json:"reviewStatus"`
}

type PaginatedFeedback struct {
	Feedbacks []Feedback `json:"feedback"`
	Metadata
}

type PostgresFeedbackStore struct {
	queries *sqlc.Queries
}

func NewPostgresFeedStore(queries *sqlc.Queries) *PostgresFeedbackStore {
	return &PostgresFeedbackStore{queries: queries}
}

type FeedbackStore interface {
	Create(context.Context, *Feedback) error
	CountRecentFeedbacksOneMonth(context.Context) (int64, error)
	GetByID(ctx context.Context, id int64) (*Feedback, error)
	List(ctx context.Context, feedbackType string, reviewStatus string, filters Filters) (*PaginatedFeedback, error)
	UpdateReviewStatus(ctx context.Context, id int64, reviewStatus string) error
}

func (pg *PostgresFeedbackStore) Create(ctx context.Context, feedback *Feedback) error {
	return pg.queries.CreateFeedback(ctx, sqlc.CreateFeedbackParams{
		Type:    feedback.Type,
		Content: feedback.Content,
	})
}

func (pg *PostgresFeedbackStore) CountRecentFeedbacksOneMonth(ctx context.Context) (int64, error) {
	return pg.queries.CountRecentFeedbacksOneMonth(ctx)
}

func (pg *PostgresFeedbackStore) GetByID(ctx context.Context, id int64) (*Feedback, error) {
	row, err := pg.queries.GetFeedbackByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Feedback{
		ID:           row.ID,
		Type:         row.Type,
		Content:      row.Content,
		CreatedAt:    row.CreatedAt.Time,
		ReviewStatus: row.ReviewStatus,
	}, nil
}

func (pg *PostgresFeedbackStore) List(ctx context.Context, feedbackType string, reviewStatus string, filters Filters) (*PaginatedFeedback, error) {

	args := sqlc.ListFeedbackParams{
		Type:         feedbackType,
		ReviewStatus: reviewStatus,
		Limit:        int32(filters.limit()),
		Offset:       int32(filters.offset()),
	}

	rows, err := pg.queries.ListFeedback(ctx, args)
	if err != nil {
		return nil, err
	}

	var totalCount = 0
	feedbackResponses := make([]Feedback, len(rows))

	for i, f := range rows {
		feedbackResponses[i] = Feedback{
			ID:           f.ID,
			Type:         f.Type,
			Content:      f.Content,
			CreatedAt:    f.CreatedAt.Time,
			ReviewStatus: f.ReviewStatus,
		}
		totalCount = int(f.Count)
	}

	return &PaginatedFeedback{
		Feedbacks: feedbackResponses,
		Metadata:  CalculateMetadata(totalCount, filters.Page, filters.PageSize),
	}, nil
}

func (pg *PostgresFeedbackStore) UpdateReviewStatus(ctx context.Context, id int64, reviewStatus string) error {
	args := sqlc.UpdateFeedbackReviewStatusParams{
		ID:           id,
		ReviewStatus: reviewStatus,
	}
	return pg.queries.UpdateFeedbackReviewStatus(ctx, args)
}

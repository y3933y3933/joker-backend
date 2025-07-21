package service

import (
	"context"
	"errors"

	"github.com/y3933y3933/joker/internal/store"
)

type FeedbackService struct {
	feedbackStore store.FeedbackStore
}

func NewFeedbackService(feedbackStore store.FeedbackStore) *FeedbackService {
	return &FeedbackService{
		feedbackStore: feedbackStore,
	}
}

func (s *FeedbackService) CreateFeedback(ctx context.Context, feedback *store.Feedback) error {
	return s.feedbackStore.Create(ctx, feedback)
}

type FeedbackQueryParams struct {
	Type         string `json:"type"`
	ReviewStatus string `json:"reviewStatus"`
	Page         int    `json:"page"`
	PageSize     int    `json:"page_size"`
}

func (s *FeedbackService) ListFeedback(ctx context.Context, query FeedbackQueryParams) (*store.PaginatedFeedback, error) {
	filters := store.Filters{
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	result, err := s.feedbackStore.List(ctx, query.Type, query.ReviewStatus, filters)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *FeedbackService) GetFeedbackByID(ctx context.Context, id int64) (*store.Feedback, error) {
	return s.feedbackStore.GetByID(ctx, id)
}

func (s *FeedbackService) UpdateFeedbackReviewStatus(ctx context.Context, id int64, reviewStatus string) error {
	return s.feedbackStore.UpdateReviewStatus(ctx, id, reviewStatus)
}

func (s *FeedbackService) ValidateFeedbackParams(params FeedbackQueryParams) error {
	// 驗證 level
	if params.Type != "" {
		if params.Type != "feature" && params.Type != "other" && params.Type != "issue" {
			return errors.New("invalid level: must be 'feature' or 'issue' or 'other'")
		}
	}

	if params.Page < 1 {
		return errors.New("page must be greater than 0")
	}

	if params.PageSize < 1 || params.PageSize > 100 {
		return errors.New("page_size must be between 1 and 100")
	}

	return nil
}

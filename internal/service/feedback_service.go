package service

import (
	"context"

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

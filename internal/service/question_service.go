package service

import (
	"context"

	"github.com/y3933y3933/joker/internal/store"
)

type QuestionService struct {
	questionStore store.QuestionStore
}

func NewQuestionService(questionStore store.QuestionStore) *QuestionService {
	return &QuestionService{
		questionStore: questionStore,
	}
}

func (s *QuestionService) ListRandomQuestions(ctx context.Context, limit int) ([]*store.Question, error) {
	return s.questionStore.ListRandomQuestions(ctx, int32(limit))
}

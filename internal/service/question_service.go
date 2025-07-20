package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/y3933y3933/joker/internal/store"
)

type QuestionQueryParams struct {
	Keyword  string `json:"keyword"`
	Level    string `json:"level"`
	SortBy   string `json:"sort_by"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

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

func (s *QuestionService) ListQuestions(ctx context.Context, query QuestionQueryParams) (*store.PaginatedQuestion, error) {

	filters := store.Filters{
		Page:     query.Page,
		PageSize: query.PageSize,
		SortBy:   query.SortBy,
	}

	result, err := s.questionStore.ListQuestions(ctx, query.Keyword, query.Level, filters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *QuestionService) ValidateQuestionParams(params QuestionQueryParams) error {
	// 驗證 level
	if params.Level != "" {
		if params.Level != "normal" && params.Level != "spicy" {
			return errors.New("invalid level: must be 'normal' or 'spicy'")
		}
	}

	validSortOptions := []string{
		"created_at_asc", "created_at_desc",
	}

	valid := false
	for _, option := range validSortOptions {
		if params.SortBy == option {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("invalid sort_by: must be one of %v", validSortOptions)
	}

	if params.Page < 1 {
		return errors.New("page must be greater than 0")
	}

	if params.PageSize < 1 || params.PageSize > 100 {
		return errors.New("page_size must be between 1 and 100")
	}

	return nil
}

func (s *QuestionService) DeleteQuestion(ctx context.Context, id int64) error {
	return s.questionStore.Delete(ctx, id)
}

func (s *QuestionService) CreateQuestion(ctx context.Context, content, level string) (*store.Question, error) {
	return s.questionStore.Create(ctx, content, level)
}

func (s *QuestionService) UpdateQuestion(ctx context.Context, id int64, content, level *string) (*store.Question, error) {
	question, err := s.questionStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if level != nil {
		question.Level = *level
	}
	if content != nil {
		question.Content = *content
	}

	return s.questionStore.Update(ctx, question.ID, question.Content, question.Level)

}

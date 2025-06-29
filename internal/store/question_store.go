package store

import (
	"context"

	"github.com/y3933y3933/joker/internal/db/sqlc"
)

type Question struct {
	ID      int64  `json:"id"`
	Level   string `json:"level"`
	Content string `json:"content"`
}

const (
	QuestionLevelNormal = "normal"
	QuestionLevelSpicy  = "spicy"
)

type PostgresQuestionStore struct {
	queries *sqlc.Queries
}

func NewPostgresQuestionStore(queries *sqlc.Queries) *PostgresQuestionStore {
	return &PostgresQuestionStore{queries: queries}
}

type QuestionStore interface {
	ListRandomQuestions(ctx context.Context, limit int32) ([]*Question, error)
}

func (pg *PostgresQuestionStore) ListRandomQuestions(ctx context.Context, limit int32) ([]*Question, error) {
	rows, err := pg.queries.ListRandomQuestions(ctx, limit)
	if err != nil {
		return nil, err
	}

	var list []*Question
	for _, r := range rows {
		list = append(list, &Question{
			ID:      r.ID,
			Level:   r.Level,
			Content: r.Content,
		})
	}
	return list, nil
}

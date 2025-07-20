package store

import (
	"context"
	"time"

	"github.com/y3933y3933/joker/internal/db/sqlc"
)

type Question struct {
	ID        int64     `json:"id"`
	Level     string    `json:"level"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type PaginatedQuestion struct {
	Questions []Question `json:"questions"`
	Metadata
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
	ListQuestions(ctx context.Context, content, level string, filters Filters) (*PaginatedQuestion, error)
	Create(ctx context.Context, content, level string) (*Question, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, content, level string) (*Question, error)
	Get(ctx context.Context, id int64) (*Question, error)
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

func (pg *PostgresQuestionStore) ListQuestions(ctx context.Context, content, level string, filters Filters) (*PaginatedQuestion, error) {
	args := sqlc.ListQuestionsParams{
		PlaintoTsquery: content,
		Level:          level,
		Column3:        filters.SortBy,
		Limit:          int32(filters.limit()),
		Offset:         int32(filters.offset()),
	}

	rows, err := pg.queries.ListQuestions(ctx, args)
	if err != nil {
		return nil, err
	}

	var totalCount = 0
	questionResponses := make([]Question, len(rows))
	for i, q := range rows {
		questionResponses[i] = Question{
			ID:        q.ID,
			Level:     q.Level,
			Content:   q.Content,
			CreatedAt: q.CreatedAt.Time,
		}
		totalCount = int(q.Count)
	}

	return &PaginatedQuestion{
		Questions: questionResponses,
		Metadata:  CalculateMetadata(totalCount, filters.Page, filters.PageSize),
	}, nil

}

func (pg *PostgresQuestionStore) Create(ctx context.Context, content, level string) (*Question, error) {
	row, err := pg.queries.CreateQuestion(ctx, sqlc.CreateQuestionParams{
		Level:   level,
		Content: content,
	})

	if err != nil {
		return nil, err
	}

	return &Question{
		ID:        row.ID,
		Content:   row.Content,
		Level:     row.Level,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (pg *PostgresQuestionStore) Delete(ctx context.Context, id int64) error {
	return pg.queries.DeleteQuestion(ctx, id)
}

func (pg *PostgresQuestionStore) Update(ctx context.Context, id int64, content, level string) (*Question, error) {
	row, err := pg.queries.UpdateQuestion(ctx, sqlc.UpdateQuestionParams{
		ID:      id,
		Level:   level,
		Content: content,
	})
	if err != nil {
		return nil, err
	}

	return &Question{
		ID:      row.ID,
		Content: row.Content,
		Level:   row.Level,
	}, nil
}

func (pg *PostgresQuestionStore) Get(ctx context.Context, id int64) (*Question, error) {
	row, err := pg.queries.GetQuestionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Question{
		ID:      row.ID,
		Level:   row.Level,
		Content: row.Content,
	}, nil
}

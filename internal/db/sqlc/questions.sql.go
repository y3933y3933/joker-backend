// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: questions.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createQuestion = `-- name: CreateQuestion :one
INSERT INTO questions (level, content)
VALUES ($1, $2)
RETURNING id, level, content, created_at, updated_at
`

type CreateQuestionParams struct {
	Level   string
	Content string
}

func (q *Queries) CreateQuestion(ctx context.Context, arg CreateQuestionParams) (Question, error) {
	row := q.db.QueryRow(ctx, createQuestion, arg.Level, arg.Content)
	var i Question
	err := row.Scan(
		&i.ID,
		&i.Level,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteQuestion = `-- name: DeleteQuestion :exec
DELETE FROM questions
WHERE id = $1
`

func (q *Queries) DeleteQuestion(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteQuestion, id)
	return err
}

const getQuestionByID = `-- name: GetQuestionByID :one
SELECT id, level, content, created_at, updated_at
FROM questions
WHERE id = $1
`

func (q *Queries) GetQuestionByID(ctx context.Context, id int64) (Question, error) {
	row := q.db.QueryRow(ctx, getQuestionByID, id)
	var i Question
	err := row.Scan(
		&i.ID,
		&i.Level,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listQuestions = `-- name: ListQuestions :many
SELECT count(*) OVER(),id, level, content, created_at
FROM questions
WHERE (to_tsvector('simple', content) @@ plainto_tsquery('simple', $1) OR $1 = '') 
AND (level = $2 OR $2 = '')
ORDER BY 
    CASE WHEN $3 = 'created_at_asc' THEN created_at END ASC,
    CASE WHEN $3 = 'created_at_desc' THEN created_at END DESC,
    id DESC
LIMIT $4 OFFSET $5
`

type ListQuestionsParams struct {
	PlaintoTsquery string
	Level          string
	Column3        interface{}
	Limit          int32
	Offset         int32
}

type ListQuestionsRow struct {
	Count     int64
	ID        int64
	Level     string
	Content   string
	CreatedAt pgtype.Timestamptz
}

func (q *Queries) ListQuestions(ctx context.Context, arg ListQuestionsParams) ([]ListQuestionsRow, error) {
	rows, err := q.db.Query(ctx, listQuestions,
		arg.PlaintoTsquery,
		arg.Level,
		arg.Column3,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListQuestionsRow
	for rows.Next() {
		var i ListQuestionsRow
		if err := rows.Scan(
			&i.Count,
			&i.ID,
			&i.Level,
			&i.Content,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRandomQuestions = `-- name: ListRandomQuestions :many
SELECT id, level, content, created_at, updated_at
FROM questions
ORDER BY RANDOM()
LIMIT $1
`

func (q *Queries) ListRandomQuestions(ctx context.Context, limit int32) ([]Question, error) {
	rows, err := q.db.Query(ctx, listRandomQuestions, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Question
	for rows.Next() {
		var i Question
		if err := rows.Scan(
			&i.ID,
			&i.Level,
			&i.Content,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateQuestion = `-- name: UpdateQuestion :one
UPDATE questions
SET level = $2, content = $3, updated_at = NOW()
WHERE id = $1
RETURNING id, level, content, created_at, updated_at
`

type UpdateQuestionParams struct {
	ID      int64
	Level   string
	Content string
}

func (q *Queries) UpdateQuestion(ctx context.Context, arg UpdateQuestionParams) (Question, error) {
	row := q.db.QueryRow(ctx, updateQuestion, arg.ID, arg.Level, arg.Content)
	var i Question
	err := row.Scan(
		&i.ID,
		&i.Level,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

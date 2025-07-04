// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: questions.sql

package sqlc

import (
	"context"
)

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

-- name: ListRandomQuestions :many
SELECT id, level, content, created_at, updated_at
FROM questions
ORDER BY RANDOM()
LIMIT $1;
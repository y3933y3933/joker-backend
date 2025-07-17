-- name: ListRandomQuestions :many
SELECT id, level, content, created_at, updated_at
FROM questions
ORDER BY RANDOM()
LIMIT $1;

-- name: ListQuestions :many
SELECT count(*) OVER(),id, level, content, created_at
FROM questions
WHERE (to_tsvector('simple', content) @@ plainto_tsquery('simple', $1) OR $1 = '') 
AND (level = $2 OR $2 = '')
ORDER BY 
    CASE WHEN $3 = 'created_at_asc' THEN created_at END ASC,
    CASE WHEN $3 = 'created_at_desc' THEN created_at END DESC,
    id DESC
LIMIT $4 OFFSET $5;


-- name: GetQuestionByID :one
SELECT id, level, content, created_at, updated_at
FROM questions
WHERE id = $1;

-- name: CreateQuestion :one
INSERT INTO questions (level, content)
VALUES ($1, $2)
RETURNING id, level, content, created_at, updated_at;

-- name: UpdateQuestion :one
UPDATE questions
SET level = $2, content = $3, updated_at = NOW()
WHERE id = $1
RETURNING id, level, content, created_at, updated_at;

-- name: DeleteQuestion :exec
DELETE FROM questions
WHERE id = $1;
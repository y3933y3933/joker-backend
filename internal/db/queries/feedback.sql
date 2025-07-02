-- name: CreateFeedback :exec
INSERT INTO feedback (type, content)
VALUES ($1, $2);
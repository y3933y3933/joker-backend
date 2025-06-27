-- name: CreateGame :one
INSERT INTO games (code, status, updated_at)
VALUES($1, $2 , updated_at = NOW())
RETURNING id, code, status, created_at;
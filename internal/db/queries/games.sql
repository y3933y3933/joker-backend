-- name: CreateGame :one
INSERT INTO games (code, status)
VALUES($1, $2)
RETURNING id, code, status, created_at;

-- name: GetGameByCode :one
SELECT id, code , status, created_at, updated_at 
FROM games
WHERE code = $1;
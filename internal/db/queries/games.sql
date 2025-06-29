-- name: CreateGame :one
INSERT INTO games (code, status)
VALUES($1, $2)
RETURNING id, code, status, created_at;

-- name: GetGameByCode :one
SELECT id, code , status, created_at, updated_at 
FROM games
WHERE code = $1;

-- name: UpdateGameStatus :exec
UPDATE games
SET status = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: EndGame :exec
UPDATE games
SET status = 'ended',
    updated_at = NOW()
WHERE code = $1;

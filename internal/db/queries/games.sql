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


-- name: GetGameStatusByID :one
SELECT status
FROM games
WHERE id = $1;

-- name: DeleteByCode :exec
DELETE FROM games WHERE code = $1;



-- name: GetGamesTodayCount :one
SELECT COUNT(*) AS games_today
FROM games
WHERE created_at >= CURRENT_DATE;



-- name: GetActiveRoomsCount :one
SELECT COUNT(*) AS active_rooms
FROM games
WHERE status != 'ended';




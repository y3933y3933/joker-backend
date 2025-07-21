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




-- name: ListGames :many
SELECT
  COUNT(*) OVER() AS total_count,
  g.id,
  g.code,
  g.status,
  COUNT(p.id) AS player_count,
  g.created_at
FROM games g
LEFT JOIN players p ON p.game_id = g.id
WHERE (UPPER(g.code) = UPPER($1) OR $1 = '')
    AND (g.status = $2 OR $2 = ''  )
GROUP BY g.id
ORDER BY g.created_at DESC
LIMIT $3 OFFSET $4;
-- name: CreatePlayer :one 
INSERT INTO players(game_id, nickname, is_host)
VALUES($1, $2, $3)
RETURNING id, game_id,nickname, is_host, joined_at;

-- name: CountPlayersInGame :one
SELECT COUNT(*)
FROM players
WHERE game_id = $1;

-- name: FindPlayersByGameID :many
SELECT id, nickname, is_host, game_id
FROM players
WHERE game_id = $1
ORDER BY id;

-- name: DeletePlayerByID :exec
DELETE FROM players WHERE id = $1;



-- name: FindPlayerByID :one
SELECT id, nickname, is_host, game_id
FROM players
WHERE id = $1;

-- name: UpdateHost :exec
UPDATE players
SET is_host = $2
WHERE id = $1;

-- name: FindPlayerByNickname :one
SELECT id, nickname, is_host, game_id
FROM players
WHERE game_id = $1 AND nickname = $2;
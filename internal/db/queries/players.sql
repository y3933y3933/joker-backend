-- name: CreatePlayer :one 
INSERT INTO players(game_id, nickname, is_host)
VALUES($1, $2, $3)
RETURNING id, game_id,nickname, is_host, joined_at;

-- name: CountPlayersInGame :one
SELECT COUNT(*)
FROM players
WHERE game_id = $1;
-- name: CreatePlayer :one 
INSERT INTO players(game_id, nickname, is_host, status)
VALUES($1, $2, $3, 'online')
RETURNING id, game_id,nickname, is_host, status;

-- name: CountPlayersInGame :one
SELECT COUNT(*)
FROM players
WHERE game_id = $1;

-- name: FindPlayersByGameID :many
SELECT id, nickname, is_host, game_id, status
FROM players
WHERE game_id = $1
ORDER BY id;

-- name: FindOnlinePlayersByGameID :many
SELECT id, nickname, game_id, is_host, status
FROM players
WHERE game_id = $1 AND status = 'online';


-- name: DeletePlayerByID :exec
DELETE FROM players WHERE id = $1;



-- name: FindPlayerByID :one
SELECT id, nickname, is_host, game_id, status
FROM players
WHERE id = $1;

-- name: UpdateHost :exec
UPDATE players
SET is_host = $2
WHERE id = $1;

-- name: FindPlayerByNickname :one
SELECT id, nickname, is_host, game_id, status
FROM players
WHERE game_id = $1 AND nickname = $2;

-- name: GetGamePlayerStats :many
SELECT
  p.id,
  p.nickname,
  COUNT(CASE WHEN r.is_joker = TRUE THEN 1 END) AS joker_cards_drawn
FROM players p
LEFT JOIN rounds r ON r.answer_player_id = p.id AND r.game_id = $1
WHERE p.game_id = $1
GROUP BY p.id, p.nickname
ORDER BY p.id;

-- name: GetPlayerCountByGameCode :one
SELECT COUNT(*) FROM players
WHERE game_id = (SELECT id FROM games WHERE code = $1);


-- name: UpdatePlayerStatus :exec
UPDATE players
SET status = $2
WHERE id = $1;

-- name: GetLivePlayerCount :one
SELECT COUNT(*) AS live_player_count
FROM players
WHERE status = 'online';
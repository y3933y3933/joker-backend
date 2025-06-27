-- name: ListPlayersByGameCode :many
SELECT p.id, p.nickname, p.is_host, p.joined_at
FROM players p
JOIN games g ON p.game_id = g.id
WHERE g.code = $1
ORDER BY p.joined_at;


-- name: CountPlayersInGame :one
SELECT COUNT(*) FROM players WHERE game_id = $1;

-- name: CreatePlayer :one
INSERT INTO players(game_id, nickname, is_host)
VALUES ($1, $2, $3)
RETURNING id, nickname, is_host, joined_at;



-- name: DeletePlayer :one
DELETE FROM players 
WHERE id = $1 AND game_id = $2
RETURNING id, nickname;


-- name: GetFirstPlayerInGame :one
SELECT id FROM players
WHERE game_id = $1
ORDER BY joined_at ASC
LIMIT 1;


-- name: GetPlayerByID :one
SELECT id, game_id, nickname, is_host, joined_at FROM players
WHERE id = $1;


-- name: GetPlayerStatsByGameID :many
SELECT
  p.id AS player_id,
  p.nickname,
  COUNT(CASE WHEN r.question_player_id = p.id THEN 1 END) AS questions_asked,
  COUNT(CASE WHEN r.answer_player_id = p.id THEN 1 END) AS questions_answered,
  COUNT(CASE WHEN r.answer_player_id = p.id AND r.is_joker = TRUE THEN 1 END) AS joker_cards_drawn
FROM players p
LEFT JOIN rounds r ON r.game_id = p.game_id
WHERE p.game_id = $1
GROUP BY p.id, p.nickname
ORDER BY p.id ASC;
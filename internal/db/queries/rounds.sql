-- name: CreateRound :one
INSERT INTO rounds (
  game_id, question_id, answer, question_player_id,
  answer_player_id, is_joker, status, deck
)
VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, game_id, question_id, answer, question_player_id,
          answer_player_id, is_joker, status, created_at, deck;
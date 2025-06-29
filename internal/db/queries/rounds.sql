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


-- name: SetRoundQuestion :exec
UPDATE rounds
SET question_id = $1,
    status = 'waiting_for_answer',
    updated_at = NOW()
WHERE id = $2;

-- name: GetRoundByID :one
SELECT id, game_id, question_id, answer, question_player_id, answer_player_id, is_joker, status, deck  
FROM rounds WHERE id = $1;

-- name: GetRoundWithQuestion :one
SELECT r.id, r.game_id, r.question_id, r.answer, r.question_player_id, r.answer_player_id, r.status, r.deck,r.is_joker, q.content AS question_content
FROM rounds r
JOIN questions q ON q.id = r.question_id
WHERE r.id = $1;

-- name: UpdateAnswer :exec
UPDATE rounds
SET answer = $2,
    status = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateDrawResult :exec
UPDATE rounds
SET is_joker = $2,
    status = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: FindLastRoundByGameID :one
SELECT id, game_id, question_id, answer, question_player_id, answer_player_id, is_joker
FROM rounds
WHERE game_id = $1
ORDER BY created_at DESC
LIMIT 1;

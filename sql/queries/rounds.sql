-- name: GetCurrentRoundByGameCode :one
SELECT r.id, r.question_player_id, r.answer_player_id, r.question_id, r.is_joker, r.created_at, g.id AS game_id,
       r.status, g.level,
       q.content AS question_content
FROM rounds r
JOIN games g ON r.game_id = g.id
JOIN questions q ON r.question_id = q.id
WHERE g.code = $1
ORDER BY r.created_at DESC
LIMIT 1;



-- name: CreateRound :one
INSERT INTO rounds (
  game_id, question_id, question_player_id, answer_player_id,deck, status
)
VALUES (
  $1, $2, $3, $4,$5, 'pending'
)
RETURNING id, question_id, question_player_id,answer_player_id, status;


-- name: GetRoundByID :one
SELECT * FROM rounds WHERE id = $1;


-- name: UpdateRoundStatus :exec
UPDATE rounds SET is_joker = $2, status = $3 WHERE id = $1;


-- name: GetLatestRoundInGame :one
SELECT * FROM rounds
WHERE game_id = $1
ORDER BY id DESC
LIMIT 1;

-- name: SetRoundQuestionID :exec
UPDATE rounds
SET question_id = $2
WHERE id = $1;


-- name: SetRoundAnswer :exec
UPDATE rounds
SET answer = $2
WHERE id = $1;

-- name: CompleteRound :exec
UPDATE rounds
SET status = $2,
    is_joker = $3
WHERE id = $1;

-- name: FinishRound :exec
UPDATE rounds
SET status = 'done'
WHERE id = $1;
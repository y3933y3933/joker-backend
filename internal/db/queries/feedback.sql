-- name: CreateFeedback :exec
INSERT INTO feedback (type, content)
VALUES ($1, $2);

-- name: CountRecentFeedbacksOneMonth :one
SELECT COUNT(*) AS feedback_count
FROM feedback
WHERE created_at >= (CURRENT_DATE - INTERVAL '30 days');


-- name: GetFeedbackByID :one
SELECT id, type, review_status ,content,created_at 
FROM feedback
WHERE id = $1;

-- name: ListFeedback :many
SELECT COUNT(*) OVER(), id, type, review_status , content, created_at 
FROM feedback 
WHERE (type = $1 OR $1 ='') AND (review_status = $2 OR $2 = '')
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;



-- name: UpdateFeedbackReviewStatus :exec
UPDATE feedback 
SET review_status = $2
WHERE id = $1;
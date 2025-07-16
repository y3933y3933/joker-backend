-- name: CreateFeedback :exec
INSERT INTO feedback (type, content)
VALUES ($1, $2);

-- name: CountRecentFeedbacksOneMonth :one
SELECT COUNT(*) AS feedback_count
FROM feedback
WHERE created_at >= (CURRENT_DATE - INTERVAL '30 days');
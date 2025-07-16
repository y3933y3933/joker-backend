-- name: CreateUser :one 
INSERT INTO users ( username, password_hash)
VALUES ($1, $2)
RETURNING id;

-- name: GetUserByUsername :one
SELECT id, username, password_hash FROM users
WHERE username = $1;

-- name: GetUserByID :one
SELECT id, username FROM users 
WHERE id = $1;
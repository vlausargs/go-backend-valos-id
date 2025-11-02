-- name: CreateUser :one
INSERT INTO users (username, email, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, username, email, password, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, username, email, password, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, username, email, password, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetAllUsers :many
SELECT id, username, email, created_at, updated_at
FROM users
ORDER BY created_at DESC;

-- name: UpdateUser :exec
UPDATE users
SET username = $2, email = $3, updated_at = $4
WHERE id = $1;

-- name: UpdatePassword :exec
UPDATE users
SET password = $2, updated_at = $3
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- name: GetUsersWithPagination :many
SELECT id, username, email, created_at, updated_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

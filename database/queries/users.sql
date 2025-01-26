-- name: CreateUser :one
INSERT INTO users (
    first_name,
    email,
    password
) VALUES ( $1, $2, $3)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUserEmail :one
UPDATE users SET
    updated_at = NOW(),
    email = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserInfo :one
UPDATE users SET
    first_name = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users SET
    password = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
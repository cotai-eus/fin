-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByKratosID :one
SELECT * FROM users WHERE kratos_identity_id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    kratos_identity_id,
    email,
    full_name,
    cpf
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
    full_name = COALESCE($2, full_name),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserBalance :exec
UPDATE users
SET balance_cents = balance_cents + $2
WHERE id = $1;

-- name: GetUserForUpdate :one
SELECT * FROM users WHERE id = $1 FOR UPDATE;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

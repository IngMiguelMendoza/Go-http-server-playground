-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: UpdateUserCredentials :one
UPDATE users
SET hashed_password = $2, email = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GeUserByMail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;


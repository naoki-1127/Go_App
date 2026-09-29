-- name: GetUserByProvider :one
SELECT * FROM identity.users
WHERE provider = $1 AND provider_user_id = $2;

-- name: CreateUser :one
INSERT INTO identity.users (name, email, provider, provider_user_id,role)
VALUES ($1, $2, $3, $4,$5)
RETURNING *;

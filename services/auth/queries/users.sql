-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUserByUUID :one
SELECT * FROM users
WHERE uuid = $1;

-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET username=$2, password=$3, role=$4
WHERE uuid = $1
RETURNING *;

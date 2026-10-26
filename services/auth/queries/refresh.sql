-- name: GetRefreshTokenByHash :one
SELECT * FROM refresh
WHERE token_hash=$1;

-- name: CreateRefreshToken :one
INSERT INTO refresh (user_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: SetRefreshTokenRevoked :one
UPDATE refresh
SET revoked_at=$2
WHERE token_hash=$1
RETURNING *;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at)
VALUES (
   $1,
   $2,
   $3
)
RETURNING token, created_at, updated_at, user_id, expires_at, revoked_at;




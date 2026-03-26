-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red;


-- name: DeleteAllUsers :exec
DELETE FROM users;


-- name: LookUpUserByEmail :one
SELECT * FROM users
WHERE email = $1;


-- name: UpdateUser :one
UPDATE users
SET email = $1,
hashed_password = $2,
updated_at = NOW()
WHERE id = $3
RETURNING *;


-- name: GetUserByID :one
SELECT id, email, created_at, updated_at, is_chirpy_red
FROM users
WHERE id = $1;


-- name: SetUserChirpyRed :exec
UPDATE users
SET is_chirpy_red = TRUE,
    updated_at = NOW()
WHERE id = $1;
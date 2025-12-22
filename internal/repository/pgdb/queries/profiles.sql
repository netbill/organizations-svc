-- name: CreateProfile :one
INSERT INTO profiles (
    account_id,
    username
) VALUES (
    sqlc.arg('account_id')::uuid,
    sqlc.arg('username')::varchar
)
RETURNING *;

-- name: UpdateProfile :one
UPDATE profiles
SET
    username = COALESCE(sqlc.narg('username')::varchar, username),
    updated_at = now()
WHERE account_id = sqlc.arg('account_id')::uuid
RETURNING *;

-- name: GetProfileByAccountID :one
SELECT *
FROM profiles
WHERE account_id = sqlc.arg('account_id')::uuid;

-- name: GetProfileByUsername :one
SELECT *
FROM profiles
WHERE username = sqlc.arg('username')::varchar;

-- name: ListProfiles :many
SELECT *
FROM profiles
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: DeleteProfile :exec
DELETE FROM profiles
WHERE account_id = sqlc.arg('account_id')::uuid;

-- name: CreateProfile :one
INSERT INTO profiles (
    account_id,
    username,
    official,
    pseudonym
) VALUES (
    sqlc.arg('account_id')::uuid,
    sqlc.arg('username')::varchar,
    sqlc.arg('official')::boolean,
    sqlc.narg('pseudonym')::varchar
)
RETURNING *;

-- name: UpdateProfile :one
UPDATE profiles
SET
    username = COALESCE(sqlc.narg('username')::varchar, username),
    official = COALESCE(sqlc.narg('official')::boolean, official),
    pseudonym = CASE
        WHEN sqlc.narg('pseudonym')::varchar IS NULL THEN pseudonym
        WHEN sqlc.narg('pseudonym')::varchar = '' THEN NULL
        ELSE sqlc.narg('pseudonym')::varchar
    END,

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

-- name: DeleteProfile :exec
DELETE FROM profiles
WHERE account_id = sqlc.arg('account_id')::uuid;

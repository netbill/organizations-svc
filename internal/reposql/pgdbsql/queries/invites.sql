-- name: CreateInvite :one
INSERT INTO invites (
    id,
    agglomeration_id,
    expires_at
) VALUES (
    sqlc.arg('id')::uuid,
    sqlc.arg('agglomeration_id')::uuid,
    sqlc.arg('expires_at')::timestamptz
)
RETURNING *;

-- name: GetInviteByID :one
SELECT *
FROM invites
WHERE id = sqlc.arg('id')::uuid;

-- name: DeleteInvite :exec
DELETE FROM invites
WHERE id = sqlc.arg('id')::uuid;

-- name: UpdateInviteStatus :one
UPDATE invites
SET
    status = sqlc.arg('status')::invite_status
WHERE id = sqlc.arg('id')::uuid
RETURNING *;



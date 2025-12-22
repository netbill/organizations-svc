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

-- name: AcceptInvite :one
UPDATE invites
SET status = 'accepted'
WHERE
    id = sqlc.arg('id')::uuid
    AND status = 'sent'
    AND expires_at > now()
RETURNING *;

-- name: DeclineInvite :one
UPDATE invites
SET status = 'declined'
WHERE
    id = sqlc.arg('id')::uuid
    AND status = 'sent'
    AND expires_at > now()
RETURNING *;



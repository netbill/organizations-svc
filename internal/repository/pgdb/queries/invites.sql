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

-- name: FilterInvitesCursor :many
SELECT *
FROM invites
WHERE
    (sqlc.narg('agglomeration_id')::uuid IS NULL OR agglomeration_id = sqlc.narg('agglomeration_id')::uuid)
    AND (sqlc.narg('status')::invite_status IS NULL OR status = sqlc.narg('status')::invite_status)
    AND (sqlc.narg('only_active')::boolean IS NULL OR (status = 'sent' AND expires_at > now()))
    AND (
        sqlc.narg('after_created_at')::timestamptz IS NULL
        OR (created_at, id) < (sqlc.narg('after_created_at')::timestamptz, sqlc.narg('after_id')::uuid)
    )
ORDER BY created_at DESC, id DESC
    LIMIT sqlc.arg('limit')::int;


-- name: CountInvites :one
SELECT COUNT(*)::bigint
FROM invites
WHERE
    (sqlc.narg('agglomeration_id')::uuid IS NULL OR agglomeration_id = sqlc.narg('agglomeration_id')::uuid)
    AND (sqlc.narg('status')::invite_status IS NULL OR status = sqlc.narg('status')::invite_status)
    AND (sqlc.narg('only_active')::boolean IS NULL OR (status = 'sent' AND expires_at > now()));

-- name: DeleteExpiredSentInvites :exec
DELETE FROM invites
WHERE status = 'sent' AND expires_at <= now();


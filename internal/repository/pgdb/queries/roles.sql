-- name: CreateRole :one
INSERT INTO roles (
    id,
    agglomeration_id,
    head,
    editable,
    rank,
    name
) VALUES (
    sqlc.arg('id')::uuid,
    sqlc.arg('agglomeration_id')::uuid,
    COALESCE(sqlc.narg('head')::boolean, false),
    COALESCE(sqlc.narg('editable')::boolean, true),
    COALESCE(sqlc.narg('rank')::int, 0),
    sqlc.arg('name')::text
)
RETURNING *;

-- name: UpdateRole :one
UPDATE roles
SET
    head = COALESCE(sqlc.narg('head')::boolean, head),
    editable = COALESCE(sqlc.narg('editable')::boolean, editable),
    rank = COALESCE(sqlc.narg('rank')::int, rank),

    name = COALESCE(sqlc.narg('name')::text, name),

    updated_at = now()
WHERE id = sqlc.arg('id')::uuid
RETURNING *;

-- name: GetRoleByID :one
SELECT *
FROM roles
WHERE id = sqlc.arg('id')::uuid;

-- name: ListRolesByAgglomerationCursor :many
SELECT *
FROM roles
WHERE
    agglomeration_id = sqlc.arg('agglomeration_id')::uuid
    AND (
        sqlc.narg('after_rank')::int IS NULL
        OR (rank, created_at, id) > (
            sqlc.narg('after_rank')::int,
            sqlc.narg('after_created_at')::timestamptz,
            sqlc.narg('after_id')::uuid
        )
    )
ORDER BY rank ASC, created_at ASC, id ASC
LIMIT sqlc.arg('limit')::int;

-- name: CountRolesByAgglomeration :one
SELECT COUNT(*)::bigint
FROM roles
WHERE agglomeration_id = sqlc.arg('agglomeration_id')::uuid;

-- name: DeleteRole :exec
DELETE FROM roles
WHERE id = sqlc.arg('id')::uuid;

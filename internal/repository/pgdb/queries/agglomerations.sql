-- name: CreateAgglomeration :one
INSERT INTO agglomerations (
    name,
    icon
) VALUES (
    sqlc.arg('name')::varchar,
    sqlc.narg('icon')::text
)
RETURNING *;

-- name: UpdateAgglomeration :one
UPDATE agglomerations
SET
    name = CASE
        WHEN sqlc.narg('name')::varchar IS NULL THEN name
        WHEN sqlc.narg('name')::varchar = '' THEN NULL
        ELSE sqlc.narg('name')::varchar
    END,

    icon = CASE
        WHEN sqlc.narg('icon')::text IS NULL THEN icon
        WHEN sqlc.narg('icon')::text = '' THEN NULL
        ELSE sqlc.narg('icon')::text
    END,

    updated_at = now()
WHERE id = sqlc.arg('id')::uuid
RETURNING *;

-- name: UpdateAgglomerationStatus :one
UPDATE agglomerations
SET
    status = sqlc.arg('status')::administration_status,
    updated_at = now()
WHERE id = sqlc.arg('id')::uuid
RETURNING *;

-- name: GetAgglomerationByID :one
SELECT *
FROM agglomerations
WHERE id = sqlc.arg('id')::uuid;

-- name: FilterAgglomerations :many
SELECT *
FROM agglomerations
WHERE
    (sqlc.narg('status')::administration_status IS NULL OR status = sqlc.narg('status')::administration_status)
    AND (sqlc.narg('name_like')::text IS NULL OR name ILIKE ('%' || sqlc.narg('name_like')::text || '%'))
    AND (
        sqlc.narg('after_created_at')::timestamptz IS NULL
        OR (created_at, id) < (sqlc.narg('after_created_at')::timestamptz, sqlc.narg('after_id')::uuid)
    )
ORDER BY created_at DESC, id DESC
    LIMIT sqlc.arg('limit')::int;

-- name: CountAgglomerations :one
SELECT COUNT(*)::bigint
FROM agglomerations
WHERE
    (sqlc.narg('status')::administration_status IS NULL OR status = sqlc.narg('status')::administration_status)
    AND (sqlc.narg('name_like')::text IS NULL OR name ILIKE ('%' || sqlc.narg('name_like')::text || '%'));

-- name: DeleteAgglomeration :exec
DELETE FROM agglomerations
WHERE id = sqlc.arg('id')::uuid;

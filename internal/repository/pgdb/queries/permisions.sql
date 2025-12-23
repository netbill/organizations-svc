-- name: GetPermissionByID :one
SELECT *
FROM permissions
WHERE id = sqlc.arg('id')::uuid;

-- name: GetPermissionByCode :one
SELECT *
FROM permissions
WHERE code = sqlc.arg('code')::varchar;

-- name: FilterPermissions :many
SELECT *
FROM permissions
WHERE
    (sqlc.narg('description')::text IS NULL OR description ILIKE (sqlc.narg('description')::text || '%'))
    AND (sqlc.narg('code')::text IS NULL OR code ILIKE (sqlc.narg('code')::text || '%'))
    AND
    (
        sqlc.narg('after_code')::varchar IS NULL
        OR (code, id) > (sqlc.narg('after_code')::varchar, sqlc.narg('after_id')::uuid)
    )
ORDER BY code ASC, id ASC
LIMIT sqlc.arg('limit')::int;

-- name: CountPermissions :one
SELECT COUNT(*)::bigint
FROM permissions
WHERE (sqlc.narg('description')::text IS NULL OR code ILIKE (sqlc.narg('description')::text || '%')) AND
      (sqlc.narg('code')::text IS NULL OR code ILIKE (sqlc.narg('code')::text || '%'));


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

-- name: UpdateRole :exec
UPDATE roles
SET
--     head = COALESCE(sqlc.narg('head')::boolean, head),
--     editable = COALESCE(sqlc.narg('editable')::boolean, editable),
    rank = COALESCE(sqlc.narg('rank')::int, rank),

    name = COALESCE(sqlc.narg('name')::text, name),

    updated_at = now()
WHERE id = sqlc.arg('id')::uuid;

-- name: GetRole :one
SELECT
    r.id,
    r.agglomeration_id,
    r.head,
    r.editable,
    r.rank,
    r.name,
    r.created_at,
    r.updated_at,

    COALESCE(
        json_agg(
            json_build_object(
                'permission_id', p.id,
                'code', p.code
            )
        ) FILTER (WHERE p.id IS NOT NULL),
        '[]'
    ) AS permissions
FROM roles r
LEFT JOIN role_permissions rp ON rp.role_id = r.id
LEFT JOIN permissions p ON p.id = rp.permission_id
WHERE r.id = $1
GROUP BY r.id;

-- name: FilterRole :many
SELECT
    r.id,
    r.agglomeration_id,
    r.head,
    r.editable,
    r.rank,
    r.name,
    r.created_at,
    r.updated_at,

    COALESCE(
        json_agg(
            json_build_object(
                'permission_id', p.id,
                'code', p.code
            )
        ) FILTER (WHERE p.id IS NOT NULL),
        '[]'
    ) AS permissions
FROM roles r
LEFT JOIN role_permissions rp ON rp.role_id = r.id
LEFT JOIN permissions p ON p.id = rp.permission_id
WHERE r.agglomeration_id = $1
    AND (
        ($2 IS NULL AND $3 IS NULL)
        OR (r.rank, r.id) > ($2, $3)
    )
GROUP BY r.id
ORDER BY r.rank ASC, r.id ASC
LIMIT $4;

-- name: DeleteRole :exec
DELETE FROM roles
WHERE id = sqlc.arg('id')::uuid;

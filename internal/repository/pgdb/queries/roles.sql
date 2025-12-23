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
    COALESCE(sqlc.arg('head')::boolean, false),
    COALESCE(sqlc.arg('editable')::boolean, true),
    COALESCE(sqlc.arg('rank')::int, 0),
    sqlc.arg('name')::text
)
RETURNING *;

-- name: UpdateRole :one
UPDATE roles
SET
--     head = COALESCE(sqlc.narg('head')::boolean, head),
--     editable = COALESCE(sqlc.narg('editable')::boolean, editable),
    rank = COALESCE(sqlc.narg('rank')::int, rank),

    name = COALESCE(sqlc.narg('name')::text, name),

    updated_at = now()
WHERE id = sqlc.arg('id')::uuid
RETURNING *;

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
                'description', p.description
            )
        ) FILTER (WHERE p.id IS NOT NULL),
        '[]'
    ) AS permissions
FROM roles r
LEFT JOIN role_permissions rp ON rp.role_id = r.id
LEFT JOIN permissions p ON p.id = rp.permission_id
WHERE r.id = $1
GROUP BY r.id;

-- name: FilterRoles :many
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
            'code', p.code,
            'description', p.description
            )
        ) FILTER (WHERE p.id IS NOT NULL),
        '[]'
    ) AS permissions
FROM roles r
LEFT JOIN role_permissions rp ON rp.role_id = r.id
LEFT JOIN permissions p ON p.id = rp.permission_id
WHERE r.agglomeration_id = sqlc.arg('agglomeration_id')::uuid

    AND (
        sqlc.narg('member_id')::uuid IS NULL
        OR EXISTS (
            SELECT 1
            FROM member_roles mr
            WHERE mr.role_id = r.id
                AND mr.member_id = sqlc.narg('member_id')::uuid
        )
    )
    AND (
        sqlc.narg('permission_codes')::text[] IS NULL
        OR EXISTS (
            SELECT 1
            FROM role_permissions rp2
            JOIN permissions p2 ON p2.id = rp2.permission_id
            WHERE rp2.role_id = r.id
                AND p2.code = ANY(sqlc.narg('permission_codes')::text[])
        )
    )

    AND (
        (sqlc.narg('cursor_rank')::int IS NULL AND sqlc.narg('cursor_id')::uuid IS NULL)
        OR (r.rank, r.id) > (sqlc.narg('cursor_rank')::int, sqlc.narg('cursor_id')::uuid)
    )

GROUP BY r.id
ORDER BY r.rank ASC, r.id ASC
LIMIT sqlc.arg('limit')::int;

-- name: CountRoles :one
SELECT COUNT(*)::bigint
FROM roles r
WHERE r.agglomeration_id = sqlc.arg('agglomeration_id')::uuid

    AND (
        sqlc.narg('member_id')::uuid IS NULL
        OR EXISTS (
            SELECT 1
            FROM member_roles mr
            WHERE mr.role_id = r.id
                AND mr.member_id = sqlc.narg('member_id')::uuid
        )
    )

    AND (
        sqlc.narg('permission_codes')::text[] IS NULL
        OR EXISTS (
            SELECT 1
            FROM role_permissions rp2
            JOIN permissions p2 ON p2.id = rp2.permission_id
            WHERE rp2.role_id = r.id
                AND p2.code = ANY(sqlc.narg('permission_codes')::text[])
    )
);


-- name: DeleteRole :exec
DELETE FROM roles
WHERE id = sqlc.arg('id')::uuid;

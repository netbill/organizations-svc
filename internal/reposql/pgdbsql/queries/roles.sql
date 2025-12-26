-- name: CreateRoleAtRank :one
WITH input AS (
    SELECT
        sqlc.arg('id')::uuid                AS id,
        sqlc.arg('agglomeration_id')::uuid  AS agglomeration_id,
        COALESCE(sqlc.arg('head')::boolean, false)     AS head,
        COALESCE(sqlc.arg('editable')::boolean, true)  AS editable,
        COALESCE(sqlc.arg('rank')::int, 0)            AS new_rank,
        sqlc.arg('name')::text              AS name
),
shift AS (
    UPDATE roles r
    SET rank = r.rank + 1,
        updated_at = now()
    WHERE r.agglomeration_id = (SELECT agglomeration_id FROM input)
        AND (SELECT new_rank FROM input) > 0
        AND r.rank >= (SELECT new_rank FROM input)
    RETURNING 1
),
ins AS (
    INSERT INTO roles (id, agglomeration_id, head, editable, rank, name)
    SELECT id, agglomeration_id, head, editable, new_rank, name
    FROM input
    RETURNING *
)
SELECT * FROM ins;

-- name: UpdateRole :one
UPDATE roles
SET
--     head = COALESCE(sqlc.narg('head')::boolean, head),
--     editable = COALESCE(sqlc.narg('editable')::boolean, editable),
    name = COALESCE(sqlc.narg('name')::text, name),

    updated_at = now()
WHERE id = sqlc.arg('id')::uuid
RETURNING *;

-- name: UpdateRoleRank :one
WITH cur AS (
    SELECT id, agglomeration_id, rank AS old_rank
    FROM roles
    WHERE id = sqlc.arg('id')::uuid
    FOR UPDATE
),
upd AS (
    UPDATE roles r
    SET
        rank = CASE
            WHEN r.id = (SELECT id FROM cur) THEN sqlc.arg('new_rank')::int

            WHEN (SELECT old_rank FROM cur) < sqlc.arg('new_rank')::int
                AND r.rank > (SELECT old_rank FROM cur)
                AND r.rank <= sqlc.arg('new_rank')::int
            THEN r.rank - 1

            WHEN (SELECT old_rank FROM cur) > sqlc.arg('new_rank')::int
                AND r.rank >= sqlc.arg('new_rank')::int
                AND r.rank < (SELECT old_rank FROM cur)
            THEN r.rank + 1

            ELSE r.rank
        END,
        updated_at = now()
    WHERE r.agglomeration_id = (SELECT agglomeration_id FROM cur)
    AND (SELECT old_rank FROM cur) <> 0
    AND sqlc.arg('new_rank')::int > 0
    RETURNING r.*
)
SELECT *
FROM upd
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

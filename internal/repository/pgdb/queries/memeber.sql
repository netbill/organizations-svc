-- name: CreateMember :one
INSERT INTO members (
    account_id,
    agglomeration_id,
    position,
    label
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetMember :one
SELECT
    m.id               AS member_id,
    m.account_id,
    m.agglomeration_id,
    m.position,
    m.label,
    m.created_at,
    m.updated_at,

    p.username,
    p.official,
    p.pseudonym,

    COALESCE(
        json_agg(
            json_build_object(
                'role_id', r.id,
                'head', r.head,
                'rank', r.rank,
                'name', r.name
            )
        ) FILTER (WHERE r.id IS NOT NULL),
    '[]'
    ) AS roles

FROM members m
JOIN profiles p ON p.account_id = m.account_id
LEFT JOIN member_roles mr ON mr.member_id = m.id
LEFT JOIN roles r ON r.id = mr.role_id
WHERE m.id = $1
GROUP BY
    m.id,
    p.account_id;

-- name: FilterMembers :many
SELECT
    m.id               AS member_id,
    m.account_id,
    m.agglomeration_id,
    m.position,
    m.label,
    m.created_at,
    m.updated_at,

    p.username,
    p.official,
    p.pseudonym,

    COALESCE(
        json_agg(
            json_build_object(
                'role_id', r.id,
                'head', r.head,
                'rank', r.rank,
                'name', r.name
            )
        ) FILTER (WHERE r.id IS NOT NULL),
        '[]'
    ) AS roles
FROM members m
JOIN profiles p ON p.account_id = m.account_id
LEFT JOIN member_roles mr ON mr.member_id = m.id
LEFT JOIN roles r ON r.id = mr.role_id
WHERE
    (sqlc.narg('agglomeration_id')::uuid IS NULL OR m.agglomeration_id = sqlc.narg('agglomeration_id')::uuid)
AND (sqlc.narg('account_id')::uuid       IS NULL OR m.account_id       = sqlc.narg('account_id')::uuid)
AND (sqlc.narg('username')::text         IS NULL OR p.username         = sqlc.narg('username')::text)

AND (
    sqlc.narg('role_id')::uuid IS NULL
    OR EXISTS (
        SELECT 1
        FROM member_roles mr2
        WHERE mr2.member_id = m.id
            AND mr2.role_id   = sqlc.narg('role_id')::uuid
    )
)

AND (
    sqlc.narg('permission_code')::text IS NULL
    OR EXISTS (
        SELECT 1
        FROM member_roles mr3
        JOIN role_permissions rp ON rp.role_id = mr3.role_id
        JOIN permissions perm ON perm.id = rp.permission_id
        WHERE mr3.member_id = m.id
            AND perm.code = sqlc.narg('permission_code')::text
    )
)

AND (
    (sqlc.narg('cursor_username')::text IS NULL AND sqlc.narg('cursor_member_id')::uuid IS NULL)
    OR (p.username, m.id) > (sqlc.narg('cursor_username')::text, sqlc.narg('cursor_member_id')::uuid)
)
GROUP BY
    m.id,
    p.account_id
ORDER BY
    p.username ASC,
    m.id ASC
LIMIT sqlc.arg('limit')::int;

-- name: CountMembers :one
SELECT COUNT(*)::bigint
FROM members m
JOIN profiles p ON p.account_id = m.account_id
WHERE
    (sqlc.narg('agglomeration_id')::uuid IS NULL OR m.agglomeration_id = sqlc.narg('agglomeration_id')::uuid)
AND (sqlc.narg('account_id')::uuid       IS NULL OR m.account_id       = sqlc.narg('account_id')::uuid)
AND (sqlc.narg('username')::text         IS NULL OR p.username         = sqlc.narg('username')::text)

AND (
    sqlc.narg('role_id')::uuid IS NULL
    OR EXISTS (
        SELECT 1
        FROM member_roles mr2
        WHERE mr2.member_id = m.id
            AND mr2.role_id   = sqlc.narg('role_id')::uuid
    )
)

AND (
    sqlc.narg('permission_code')::text IS NULL
    OR EXISTS (
        SELECT 1
        FROM member_roles mr3
        JOIN role_permissions rp ON rp.role_id = mr3.role_id
        JOIN permissions perm ON perm.id = rp.permission_id
        WHERE mr3.member_id = m.id
            AND perm.code = sqlc.narg('permission_code')::text
    )
);


-- name: MemberExists :one
SELECT EXISTS (
    SELECT 1
    FROM members
    WHERE account_id = $1
      AND agglomeration_id = $2
);

-- name: UpdateMember :one
UPDATE members
SET
    position = COALESCE($2, position),
    label = COALESCE($3, label),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteMember :exec
DELETE FROM members
WHERE id = $1;
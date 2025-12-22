-- name: CreateMember :exec
INSERT INTO members (
    account_id,
    agglomeration_id,
    position,
    label
) VALUES (
    $1, $2, $3, $4
);

-- name: MemberExists :one
SELECT EXISTS (
    SELECT 1
    FROM members
    WHERE account_id = $1
      AND agglomeration_id = $2
);

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

-- name: ListMembers :many
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
WHERE m.agglomeration_id = $1
    AND (
        ($2 IS NULL AND $3 IS NULL)
        OR (p.username, m.id) > ($2, $3)
    )
GROUP BY
    m.id,
    p.account_id
ORDER BY
    p.username ASC,
    m.id ASC
LIMIT $4;

-- name: CreateOrgRoleRankShift :many
WITH bump AS (
	UPDATE organization_role_ranks
	SET rank = rank + 1000000,
		updated_at = (now() AT TIME ZONE 'UTC')
	WHERE organization_id = $1
		AND rank >= $3
),
shift AS (
    UPDATE organization_role_ranks
    SET rank = rank - 999999,
        updated_at = (now() AT TIME ZONE 'UTC')
    WHERE organization_id = $1
    	AND rank >= 1000000 + $3
),
ins AS (
    INSERT INTO organization_role_ranks (organization_id, role_id, rank)
    VALUES ($1, $2, $3)
    RETURNING organization_id, role_id, rank, created_at, updated_at
)
SELECT organization_id, role_id, rank, created_at, updated_at
FROM organization_role_ranks
WHERE organization_id = $1
ORDER BY rank ASC;

-- name: GetOrgRoleRank :one
SELECT organization_id, role_id, rank, created_at, updated_at
FROM organization_role_ranks
WHERE role_id = $1
LIMIT 1;

-- name: GetOrgRolesRanks :many
SELECT organization_id, role_id, rank, created_at, updated_at
FROM organization_role_ranks
WHERE organization_id = $1
ORDER BY rank ASC;

-- name: UpdateOrgRolesRanks :many
WITH input AS (
	SELECT *
	FROM unnest($2::uuid[], $3::int4[]) AS t(role_id, rank)
),
bump AS (
	UPDATE organization_role_ranks rr
	SET rank = rr.rank + 1000000,
		updated_at = (now() AT TIME ZONE 'UTC')
	FROM input i
	WHERE rr.organization_id = $1
		AND rr.role_id = i.role_id
),
setr AS (
	UPDATE organization_role_ranks rr
	SET rank = i.rank,
		updated_at = (now() AT TIME ZONE 'UTC')
	FROM input i
	WHERE rr.organization_id = $1
		AND rr.role_id = i.role_id
)
SELECT organization_id, role_id, rank, created_at, updated_at
FROM organization_role_ranks
WHERE organization_id = $1
ORDER BY rank ASC;


-- name: DeleteOrgRoleRankShift :many
WITH target AS (
    SELECT organization_id, role_id, rank
    FROM organization_role_ranks
    WHERE organization_id = $1
		AND role_id = $2
    FOR UPDATE
),
del AS (
    DELETE FROM organization_role_ranks rr
    USING target t
    WHERE rr.organization_id = t.organization_id
		AND rr.role_id = t.role_id
    RETURNING 1
),
bump AS (
    UPDATE organization_role_ranks rr
    SET rank = rr.rank + 1000000,
        updated_at = (now() AT TIME ZONE 'UTC')
    FROM target t
    WHERE rr.organization_id = t.organization_id
    	AND rr.rank > t.rank
),
shift AS (
    UPDATE organization_role_ranks rr
    SET rank = rr.rank - 1000001,
        updated_at = (now() AT TIME ZONE 'UTC')
    FROM target t
    WHERE rr.organization_id = t.organization_id
    	AND rr.rank > 1000000 + t.rank
)
SELECT organization_id, role_id, rank, created_at, updated_at
FROM organization_role_ranks
WHERE organization_id = $1
ORDER BY rank ASC;
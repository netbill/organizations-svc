-- name: CreateCity :one
INSERT INTO cities (
    agglomeration_id,
    slug,
    name,
    icon,
    banner
) VALUES (
    sqlc.narg('agglomeration_id')::uuid,
    sqlc.narg('slug')::varchar,
    sqlc.arg('name')::varchar,
    sqlc.narg('icon')::text,
    sqlc.narg('banner')::text
)
RETURNING *;

-- name: UpdateCity :one
UPDATE cities
SET
    agglomeration_id = COALESCE(sqlc.narg('agglomeration_id')::uuid, agglomeration_id),

    name = COALESCE(sqlc.narg('name')::varchar, name),

    slug = CASE
        WHEN sqlc.narg('slug')::varchar IS NULL THEN slug
        WHEN sqlc.narg('slug')::varchar = '' THEN NULL
        ELSE sqlc.narg('slug')::varchar
    END,

    icon = CASE
        WHEN sqlc.narg('icon')::text IS NULL THEN icon
        WHEN sqlc.narg('icon')::text = '' THEN NULL
        ELSE sqlc.narg('icon')::text
    END,

    banner = CASE
        WHEN sqlc.narg('banner')::text IS NULL THEN banner
        WHEN sqlc.narg('banner')::text = '' THEN NULL
        ELSE sqlc.narg('banner')::text
    END,

    updated_at = now()
WHERE id = sqlc.arg('id')::uuid
RETURNING *;

-- name: ActivateCity :one
UPDATE cities
SET
    status = 'active',
    updated_at = now()
WHERE id = sqlc.arg('id')::uuid
RETURNING *;

-- name: DeactivateCity :one
UPDATE cities
SET
    status = 'inactive',
    updated_at = now()
WHERE id = sqlc.arg('id')::uuid
RETURNING *;

-- name: GetCityByID :one
SELECT *
FROM cities
WHERE id = sqlc.arg('id')::uuid;

-- name: GetCityBySlug :one
SELECT *
FROM cities
WHERE slug = sqlc.arg('slug')::varchar;

-- name: FilterCities :many
SELECT *
FROM cities
WHERE
    (sqlc.narg('agglomeration_id')::uuid IS NULL OR agglomeration_id = sqlc.narg('agglomeration_id')::uuid)
    AND (sqlc.narg('status')::cities_status IS NULL OR status = sqlc.narg('status')::cities_status)
    AND (sqlc.narg('name_like')::text IS NULL OR name ILIKE ('%' || sqlc.narg('name_like')::text || '%'))
    AND (
        sqlc.narg('after_created_at')::timestamptz IS NULL
        OR (created_at, id) < (sqlc.narg('after_created_at')::timestamptz, sqlc.narg('after_id')::uuid)
    )
ORDER BY created_at DESC, id DESC
    LIMIT sqlc.arg('limit')::int;


-- name: CountCities :one
SELECT COUNT(*)::bigint
FROM cities
WHERE
    (sqlc.narg('agglomeration_id')::uuid IS NULL OR agglomeration_id = sqlc.narg('agglomeration_id')::uuid)
    AND (sqlc.narg('status')::cities_status IS NULL OR status = sqlc.narg('status')::cities_status)
    AND (sqlc.narg('name_like')::text IS NULL OR name ILIKE ('%' || sqlc.narg('name_like')::text || '%'));

-- name: DeleteCity :exec
DELETE FROM cities
WHERE id = sqlc.arg('id')::uuid;

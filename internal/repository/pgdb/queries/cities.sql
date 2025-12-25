-- name: CreateCity :exec
INSERT INTO cities (
    agglomeration_id,
    slug,
    name,
    icon,
    banner,
    point
) VALUES (
    sqlc.narg('agglomeration_id')::uuid,
    sqlc.narg('slug')::varchar,
    sqlc.arg('name')::varchar,
    sqlc.narg('icon')::text,
    sqlc.narg('banner')::text,

    ST_SetSRID(
        ST_MakePoint(
            sqlc.arg('point_lng')::double precision,
            sqlc.arg('point_lat')::double precision
        ),
        4326
    )::geography
);

-- name: UpdateCity :exec
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

    point = CASE
        WHEN sqlc.narg('point_lng')::double precision IS NULL OR sqlc.narg('point_lat')::double precision IS NULL
            THEN point
        ELSE ST_SetSRID(ST_MakePoint(sqlc.narg('point_lng')::double precision, sqlc.narg('point_lat')::double precision), 4326)::geography
    END,

    updated_at = now()
WHERE id = sqlc.arg('id')::uuid;

-- name: UpdateCityStatus :exec
UPDATE cities
SET
    status = sqlc.arg('status')::cities_status,
    updated_at = now()
WHERE id = sqlc.arg('id')::uuid;

-- name: GetCityByID :one
SELECT
    c.id,
    c.agglomeration_id,
    c.status,
    c.slug,
    c.name,
    c.icon,
    c.banner,

    ST_X(c.point::geometry)::double precision AS point_lng,
    ST_Y(c.point::geometry)::double precision AS point_lat,

    c.created_at,
    c.updated_at
FROM cities c
WHERE c.id = sqlc.arg('id')::uuid;

-- name: GetCityBySlug :one
SELECT
    c.id,
    c.agglomeration_id,
    c.status,
    c.slug,
    c.name,
    c.icon,
    c.banner,

    ST_X(c.point::geometry)::double precision AS point_lng,
    ST_Y(c.point::geometry)::double precision AS point_lat,

    c.created_at,
    c.updated_at
FROM cities c
WHERE c.slug = sqlc.arg('slug')::varchar;

-- name: FilterCities :many
SELECT
    c.id,
    c.agglomeration_id,
    c.status,
    c.slug,
    c.name,
    c.icon,
    c.banner,

    ST_X(c.point::geometry)::double precision AS point_lng,
    ST_Y(c.point::geometry)::double precision AS point_lat,

    c.created_at,
    c.updated_at
FROM cities c
WHERE
    (sqlc.narg('agglomeration_id')::uuid IS NULL OR c.agglomeration_id = sqlc.narg('agglomeration_id')::uuid)
    AND (sqlc.narg('status')::cities_status IS NULL OR c.status = sqlc.narg('status')::cities_status)
    AND (sqlc.narg('name_like')::text IS NULL OR c.name ILIKE ('%' || sqlc.narg('name_like')::text || '%'))
    AND (
        sqlc.narg('after_created_at')::timestamptz IS NULL
        OR (c.created_at, c.id) < (sqlc.narg('after_created_at')::timestamptz, sqlc.narg('after_id')::uuid)
    )
ORDER BY c.created_at DESC, c.id DESC
    LIMIT sqlc.arg('limit')::int;

-- name: FilterCitiesNearest :many
WITH u AS (
    SELECT ST_SetSRID(
        ST_MakePoint(
            sqlc.arg('user_lng')::double precision,
            sqlc.arg('user_lat')::double precision
        ),
        4326
    )::geography AS up
),
    base AS (
        SELECT
            c.id,
            c.agglomeration_id,
            c.status,
            c.slug,
            c.name,
            c.icon,
            c.banner,

            ST_X(c.point::geometry)::double precision AS point_lng,
            ST_Y(c.point::geometry)::double precision AS point_lat,

            c.created_at,
            c.updated_at,

            floor(ST_Distance(c.point, u.up))::bigint AS distance_m
        FROM cities c
        CROSS JOIN u
        WHERE
            (sqlc.narg('agglomeration_id')::uuid IS NULL OR c.agglomeration_id = sqlc.narg('agglomeration_id')::uuid)
            AND (sqlc.narg('status')::cities_status IS NULL OR c.status = sqlc.narg('status')::cities_status)
            AND (sqlc.narg('name_like')::text IS NULL OR c.name ILIKE ('%' || sqlc.narg('name_like')::text || '%'))
)
SELECT *
FROM base
WHERE
    (
        sqlc.narg('after_distance_m')::bigint IS NULL
        OR (
            distance_m > sqlc.narg('after_distance_m')::bigint
            OR (
                distance_m = sqlc.narg('after_distance_m')::bigint
                AND sqlc.narg('after_id')::uuid IS NOT NULL
                AND id > sqlc.narg('after_id')::uuid
            )
        )
    )
ORDER BY distance_m ASC, id ASC
LIMIT sqlc.arg('limit')::int;

-- name: CountCities :one
SELECT COUNT(*)::bigint
FROM cities
WHERE
    (sqlc.narg('agglomeration_id')::uuid IS NULL OR agglomeration_id = sqlc.narg('agglomeration_id')::uuid)
    AND (sqlc.narg('status')::cities_status IS NULL OR status = sqlc.narg('status')::cities_status)
    AND (sqlc.narg('name_like')::text IS NULL OR name ILIKE ('%' || sqlc.narg('name_like')::text || '%'));

-- name: CountCitiesNearest :one
WITH u AS (
    SELECT ST_SetSRID(
        ST_MakePoint(
            sqlc.arg('user_lng')::double precision,
            sqlc.arg('user_lat')::double precision
        ),
    4326
    )::geography AS up
)
SELECT COUNT(*)::bigint
FROM cities c
CROSS JOIN u
WHERE
    c.point IS NOT NULL
    AND (sqlc.narg('agglomeration_id')::uuid IS NULL OR c.agglomeration_id = sqlc.narg('agglomeration_id')::uuid)
    AND (sqlc.narg('status')::cities_status IS NULL OR c.status = sqlc.narg('status')::cities_status)
    AND (sqlc.narg('name_like')::text IS NULL OR c.name ILIKE ('%' || sqlc.narg('name_like')::text || '%'));
--     AND (
--         sqlc.narg('radius_m')::double precision IS NULL
--         OR ST_DWithin(c.point, u.up, sqlc.narg('radius_m')::double precision)
--     );

-- name: DeleteCity :exec
DELETE FROM cities
WHERE id = sqlc.arg('id')::uuid;

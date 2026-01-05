-- +migrate Up
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE place_classes (
    id           UUID                 PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_id    UUID                 REFERENCES place_classes(id) ON DELETE SET NULL,
    icon         VARCHAR(255)         NOT NULL,
    name         VARCHAR(255)         NOT NULL,

    created_at   TIMESTAMPTZ   NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at   TIMESTAMPTZ   NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),

    UNIQUE (parent_id, name),
    CHECK (parent_id IS NULL OR parent_id <> id)
);

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION trg_place_classes_no_cycles()
RETURNS trigger AS $$
DECLARE
    cycle_exists boolean;
BEGIN
    IF NEW.parent_id IS NULL THEN
        RETURN NEW;
    END IF;

    IF TG_OP = 'UPDATE' AND NEW.parent_id = OLD.parent_id THEN
        RETURN NEW;
    END IF;

    IF NEW.parent_id = NEW.id THEN
        RAISE EXCEPTION 'place_classes cycle: parent_id cannot equal id';
    END IF;

    WITH RECURSIVE ancestors AS (
        SELECT pc.id, pc.parent_id
        FROM place_classes pc
        WHERE pc.id = NEW.parent_id

        UNION ALL

        SELECT pc2.id, pc2.parent_id
        FROM place_classes pc2
        JOIN ancestors a ON a.parent_id = pc2.id
    )
    SELECT EXISTS (SELECT 1 FROM ancestors WHERE id = NEW.id)
    INTO cycle_exists;

    IF cycle_exists THEN
        RAISE EXCEPTION 'place_classes cycle detected for class %', NEW.id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER place_classes_no_cycles
BEFORE INSERT OR UPDATE OF parent_id ON place_classes
FOR EACH ROW
EXECUTE FUNCTION trg_place_classes_no_cycles();
-- +migrate StatementEnd

CREATE type "place_statuses" AS ENUM (
    'active',
    'inactive',
    'suspended'
);

CREATE TABLE places (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    class_id        UUID NOT NULL REFERENCES place_classes(id) ON DELETE RESTRICT ON UPDATE CASCADE,

    status     place_statuses         NOT NULL,
    verified   BOOLEAN                NOT NULL DEFAULT FALSE,
    point      geography(POINT, 4326) NOT NULL,

    name        VARCHAR(128) NOT NULL,
    address     VARCHAR(255) NOT NULL,
    description VARCHAR(1024),
    icon        VARCHAR(255),
    banner      VARCHAR(255),
    website     VARCHAR(255),
    phone       VARCHAR(32),

    created_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE TABLE place_features (
    id          UUID          PRIMARY KEY,
    code        VARCHAR(255)  UNIQUE NOT NULL,
    description VARCHAR(1024) NOT NULL
);

CREATE TABLE place_feature_links (
     place_id   UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
     feature_id UUID NOT NULL REFERENCES place_features(id) ON DELETE CASCADE,
     PRIMARY KEY (place_id, feature_id)
);

CREATE TABLE place_timetables (
    id        UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    place_id  UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    start_min INT  NOT NULL,
    end_min   INT  NOT NULL,

    CHECK (start_min >= 0 AND end_min <= 10080 AND end_min > start_min),

    EXCLUDE USING gist (
        place_id WITH =,
        int4range(start_min, end_min, '[)') WITH &&
    )
);

-- +migrate Down
DROP TRIGGER IF EXISTS place_classes_no_cycles ON place_classes;

DROP FUNCTION IF EXISTS trg_place_classes_no_cycles();

DROP TABLE IF EXISTS place_feature_links CASCADE;
DROP TABLE IF EXISTS place_features CASCADE;
DROP TABLE IF EXISTS place_timetables CASCADE;
DROP TABLE IF EXISTS places CASCADE;
DROP TABLE IF EXISTS place_classes CASCADE;

DROP TYPE IF EXISTS place_statuses;

DROP EXTENSION IF EXISTS btree_gist;

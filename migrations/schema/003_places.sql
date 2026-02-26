-- +migrate Up

CREATE TYPE place_statuses AS ENUM (
    'active',
    'inactive',
    'suspended'
);

CREATE TABLE places (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    class_id        UUID NOT NULL REFERENCES place_classes(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    status   place_statuses         NOT NULL DEFAULT 'inactive',
    verified BOOLEAN                NOT NULL DEFAULT FALSE,
    point    geography(POINT, 4326) NOT NULL,
    address  VARCHAR(255)           NOT NULL,
    name     VARCHAR(128)           NOT NULL,
    version  INT                    NOT NULL DEFAULT 1 CHECK (version > 0),

    description VARCHAR(1024),
    icon_key    VARCHAR(255),
    banner_key  VARCHAR(255),
    website     VARCHAR(255),
    phone       VARCHAR(32),

    source_created_at  TIMESTAMPTZ NOT NULL,
    source_updated_at  TIMESTAMPTZ NOT NULL,
    replica_created_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    replica_updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc')
);
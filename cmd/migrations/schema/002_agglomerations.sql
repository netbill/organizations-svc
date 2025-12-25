-- +migrate Up
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TYPE administration_status AS ENUM (
    'active',
    'inactive',
    'suspended'
);

CREATE TABLE agglomerations (
    id         UUID                  PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    status     administration_status NOT NULL DEFAULT 'active',
    name       VARCHAR(255)          NOT NULL,
    icon       TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TYPE cities_status AS ENUM (
    'active',
    'inactive',
    'archived'
);

CREATE TABLE cities (
    id               UUID          PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    agglomeration_id UUID          REFERENCES agglomerations(id) ON DELETE SET NULL,
    status           cities_status NOT NULL DEFAULT 'active',
    slug             VARCHAR(255)  UNIQUE,
    name             VARCHAR(255)  NOT NULL,
    icon             TEXT,
    banner           TEXT,

    point geography(Point, 4326) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE profiles (
    account_id  UUID        PRIMARY KEY,
    username    VARCHAR(32) NOT NULL UNIQUE,
    official    BOOLEAN NOT NULL DEFAULT FALSE,
    pseudonym   VARCHAR(128),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc')
);

CREATE TABLE members (
    id               UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    account_id       UUID NOT NULL REFERENCES profiles(account_id) ON DELETE CASCADE,
    agglomeration_id UUID NOT NULL REFERENCES agglomerations(id) ON DELETE CASCADE,
    position         TEXT,
    label            TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),

    UNIQUE(account_id, agglomeration_id)
);

CREATE TABLE roles (
    id               UUID    PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    agglomeration_id UUID    NOT NULL REFERENCES agglomerations(id) ON DELETE CASCADE,
    head             BOOLEAN NOT NULL DEFAULT false,
    editable         BOOLEAN NOT NULL DEFAULT true,
    rank             INT     NOT NULL DEFAULT 0 CHECK ( rank >= 0 ),
    name             TEXT    NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(agglomeration_id, name)
);

ALTER TABLE roles
    ADD CONSTRAINT roles_agglomeration_id_rank_key
    UNIQUE (agglomeration_id, rank)
    DEFERRABLE INITIALLY DEFERRED;

CREATE UNIQUE INDEX roles_one_head_per_agglomeration
    ON roles (agglomeration_id)
    WHERE head = true;

CREATE TABLE member_roles (
    member_id UUID NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    role_id   UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,

    PRIMARY KEY (member_id, role_id)
);

CREATE TABLE permissions (
    id          UUID          PRIMARY KEY,
    code        VARCHAR(255)  UNIQUE NOT NULL,
    description VARCHAR(1024) NOT NULL
);

CREATE TABLE role_permissions (
    role_id       UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,

    PRIMARY KEY (role_id, permission_id)
);

CREATE TYPE invite_status AS ENUM (
    'sent',
    'declined',
    'accepted'
);

CREATE TABLE invites (
    id               UUID          PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    agglomeration_id UUID          NOT NULL REFERENCES agglomerations(id) ON DELETE CASCADE,
    status           invite_status NOT NULL DEFAULT 'sent',

    expires_at       TIMESTAMPTZ NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc')
);

-- +migrate Down
DROP TABLE IF EXISTS agglomerations CASCADE;
DROP TABLE IF EXISTS cities CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS profiles CASCADE;
DROP TABLE IF EXISTS members CASCADE;
DROP TABLE IF EXISTS member_roles CASCADE;
DROP TABLE IF EXISTS invites CASCADE;
DROP TABLE IF EXISTS permissions CASCADE;
DROP TABLE IF EXISTS role_permissions CASCADE;

DROP TYPE IF EXISTS administration_status;
DROP TYPE IF EXISTS cities_status;
DROP TYPE IF EXISTS invite_status;
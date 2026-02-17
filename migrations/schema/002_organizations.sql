-- +migrate Up
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE profiles (
    account_id  UUID        PRIMARY KEY,
    username    VARCHAR(32) NOT NULL UNIQUE,
    official    BOOLEAN NOT NULL DEFAULT FALSE,
    pseudonym   VARCHAR(128),
    avatar_key  TEXT,

    source_created_at  TIMESTAMPTZ NOT NULL,
    source_updated_at  TIMESTAMPTZ NOT NULL,
    replica_created_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    replica_updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc')
);

CREATE TYPE organization_status AS ENUM (
    'active',
    'inactive',
);

CREATE TABLE organizations (
    id         UUID                  PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    status     organization_status   NOT NULL DEFAULT 'active',
    name       VARCHAR(255)          NOT NULL,
    icon_key   TEXT,
    banner_key TEXT,
    max_roles  INT                   NOT NULL DEFAULT 100 CHECK ( max_roles > 0 ),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE organization_members (
    id               UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    account_id       UUID NOT NULL REFERENCES profiles(account_id) ON DELETE CASCADE,
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    head             BOOLEAN NOT NULL DEFAULT FALSE,
    position         VARCHAR(255),
    label            VARCHAR(128),

    created_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),

    UNIQUE(account_id, organization_id)
);

CREATE UNIQUE INDEX members_one_head_per_organization
    ON organization_members (organization_id)
    WHERE head = true;

CREATE TYPE organization_invite_status AS ENUM (
    'sent',
    'declined',
    'accepted'
);

CREATE TABLE organization_invites (
    id               UUID          PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    account_id       UUID          NOT NULL REFERENCES profiles(account_id) ON DELETE CASCADE,
    organization_id  UUID          NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    status           organization_invite_status NOT NULL DEFAULT 'sent',

    expires_at       TIMESTAMPTZ NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc')
);

-- +migrate Down
DROP TABLE IF EXISTS organization_members CASCADE;
DROP TABLE IF EXISTS organization_invites CASCADE;
DROP TABLE IF EXISTS organizations CASCADE;
DROP TABLE IF EXISTS profiles CASCADE;

DROP INDEX IF EXISTS members_one_head_per_organization;

DROP TYPE IF EXISTS organization_status;
DROP TYPE IF EXISTS organization_invite_status;
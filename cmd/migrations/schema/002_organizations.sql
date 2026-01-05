-- +migrate Up
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE profiles (
    account_id  UUID        PRIMARY KEY,
    username    VARCHAR(32) NOT NULL UNIQUE,
    official    BOOLEAN NOT NULL DEFAULT FALSE,
    pseudonym   VARCHAR(128),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc')
);

CREATE TYPE organization_status AS ENUM (
    'active',
    'inactive',
    'suspended'
);

CREATE TABLE organizations (
    id         UUID                  PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    status     organization_status   NOT NULL DEFAULT 'active',
    verified   BOOLEAN               NOT NULL DEFAULT FALSE,
    name       VARCHAR(255)          NOT NULL,
    icon       TEXT,
    max_roles  INT                   NOT NULL DEFAULT 100 CHECK ( max_roles > 0 ),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE organization_members (
    id               UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    account_id       UUID NOT NULL REFERENCES profiles(account_id) ON DELETE CASCADE,
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    position         TEXT,
    label            TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),

    UNIQUE(account_id, organization_id)
);

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

DROP TYPE IF EXISTS organization_status;
DROP TYPE IF EXISTS organization_invite_status;
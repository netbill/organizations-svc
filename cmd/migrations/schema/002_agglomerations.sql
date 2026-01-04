-- +migrate Up
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE administration_status AS ENUM (
    'active',
    'inactive',
    'suspended'
);

CREATE TABLE organizations (
    id         UUID                  PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    status     administration_status NOT NULL DEFAULT 'active',
    verified   BOOLEAN               NOT NULL DEFAULT FALSE,
    name       VARCHAR(255)          NOT NULL,
    icon       TEXT,
    max_roles  INT                   NOT NULL DEFAULT 100 CHECK ( max_roles > 0 ),

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
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    position         TEXT,
    label            TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),

    UNIQUE(account_id, organization_id)
);

CREATE TABLE roles (
    id               UUID    PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    organization_id  UUID    NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    head             BOOLEAN NOT NULL DEFAULT false,
    rank             INT     NOT NULL DEFAULT 0 CHECK ( rank >= 0 ),
    name             TEXT    NOT NULL,
    description      TEXT    NOT NULL,
    color            TEXT    NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(organization_id, name)
);

CREATE UNIQUE INDEX roles_one_head_per_organization
    ON roles (organization_id)
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

INSERT INTO permissions (id, code, description) VALUES
    (uuid_generate_v4(), 'organization.manage', 'manage organization settings'),
    (uuid_generate_v4(), 'invites.manage', 'manage organization invites'),
    (uuid_generate_v4(), 'members.manage', 'manage organization members'),
    (uuid_generate_v4(), 'roles.manage', 'manage organization roles');


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
    account_id       UUID          NOT NULL REFERENCES profiles(account_id) ON DELETE CASCADE,
    organization_id  UUID          NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    status           invite_status NOT NULL DEFAULT 'sent',

    expires_at       TIMESTAMPTZ NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc')
);

-- +migrate Down
DROP TABLE IF EXISTS organizations CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS profiles CASCADE;
DROP TABLE IF EXISTS members CASCADE;
DROP TABLE IF EXISTS member_roles CASCADE;
DROP TABLE IF EXISTS invites CASCADE;
DROP TABLE IF EXISTS permissions CASCADE;
DROP TABLE IF EXISTS role_permissions CASCADE;

DROP TYPE IF EXISTS administration_status;
DROP TYPE IF EXISTS invite_status;
-- +migrate Up
CREATE TABLE organization_roles (
    id              UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    name        TEXT NOT NULL,
    description TEXT NOT NULL,
    color       TEXT NOT NULL,
    version     INT  NOT NULL DEFAULT 1 CHECK (version > 0),

    created_at  TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),

    UNIQUE (organization_id, name)
);

CREATE TABLE organization_role_ranks (
    organization_id UUID NOT NULL,
    role_id         UUID NOT NULL,
    rank            INT  NOT NULL CHECK (rank >= 0),

    created_at      TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),

    PRIMARY KEY (role_id),
    UNIQUE (organization_id, rank),

    FOREIGN KEY (organization_id, role_id)
        REFERENCES organization_roles (organization_id, id)
        ON DELETE CASCADE
);

CREATE TABLE organization_role_ranks_revisions (
    organization_id UUID PRIMARY KEY REFERENCES organizations(id) ON DELETE CASCADE,
    revision        INT   NOT NULL DEFAULT 0,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE TABLE organization_member_role_links (
    member_id  UUID NOT NULL REFERENCES organization_members(id) ON DELETE CASCADE,
    role_id    UUID NOT NULL REFERENCES organization_roles (id) ON DELETE CASCADE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    PRIMARY KEY (member_id, role_id)
);

CREATE TABLE organization_member_role_links_revisions (
    member_id  UUID PRIMARY KEY REFERENCES organization_members(id) ON DELETE CASCADE,
    revision   INT NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE TABLE organization_role_permissions ( -- тут все гуд
    id          UUID          PRIMARY KEY NOT NULL,
    code        VARCHAR(255)  NOT NULL UNIQUE,
    description VARCHAR(1024) NOT NULL,

    created_at  TIMESTAMPTZ   NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

INSERT INTO organization_role_permissions (id, code, description) VALUES
    ('8a1f6a63-2a5d-4b35-9f64-2f2c3bce9e01', 'organization.update', 'manage organization settings'),
    ('7d2c4e91-1c5a-4b23-a9b1-93df0a12f402', 'roles.manage',        'manage organization roles'),
    ('c1e0a572-9a2f-4e9c-b67e-5d01a7c44d03', 'invites.manage',      'manage organization invites'),
    ('2b1d5c64-6f12-4c78-9c34-8b2e6d0f1a04', 'members.delete',      'remove organization members'),
    ('a4e7c903-3d11-4b8f-8e1d-12f6b2d9c205', 'members.update',      'update organization members'),
    ('d6f3a890-7b2c-4e21-b93a-56a1c8d0e306', 'places.create',       'create places within the organization'),
    ('e7a4b901-8c3d-4f12-a92b-67b2d9e1f407', 'places.delete',       'delete places within the organization'),
    ('f8b5c012-9d4e-4a23-b81c-78c3e0f2a508', 'places.update',       'update places within the organization')
ON CONFLICT (id) DO UPDATE
SET code = EXCLUDED.code,
    description = EXCLUDED.description,
    updated_at = (now() AT TIME ZONE 'UTC');

CREATE TABLE organization_role_permission_links (
    role_id       UUID NOT NULL REFERENCES organization_roles (id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES organization_role_permissions (id) ON DELETE RESTRICT,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE organization_role_permission_links_revisions (
    role_id    UUID PRIMARY KEY REFERENCES organization_roles(id) ON DELETE CASCADE,
    revision   INT NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

-- +migrate Down
DROP TABLE IF EXISTS organization_role_permission_links_revisions;
DROP TABLE IF EXISTS organization_role_permission_links;
DROP TABLE IF EXISTS organization_role_permissions;

DROP TABLE IF EXISTS organization_member_role_links_revisions;
DROP TABLE IF EXISTS organization_member_role_links;

DROP TABLE IF EXISTS organization_role_ranks_revisions;
DROP TABLE IF EXISTS organization_role_ranks;

DROP TABLE IF EXISTS organization_roles;
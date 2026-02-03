-- +migrate Up
CREATE TABLE organization_roles (
    id               UUID    PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    organization_id  UUID    NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    rank             INT     NOT NULL DEFAULT 0 CHECK (rank >= 0),
    name             TEXT    NOT NULL,
    description      TEXT    NOT NULL,
    color            TEXT    NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (organization_id, name)
);

CREATE TABLE organization_member_roles (
    member_id UUID NOT NULL REFERENCES organization_members(id) ON DELETE CASCADE,
    role_id   UUID NOT NULL REFERENCES organization_roles (id) ON DELETE CASCADE,

    PRIMARY KEY (member_id, role_id)
);

CREATE TABLE organization_role_permissions (
    code        VARCHAR(255)  PRIMARY KEY UNIQUE NOT NULL,
    description VARCHAR(1024) NOT NULL
);

INSERT INTO organization_role_permissions (code, description) VALUES
    ('organization.manage', 'manage organization settings'),
    ('invites.manage', 'manage organization invites'),
    ('members.manage', 'manage organization members'),
    ('roles.manage', 'manage organization roles');

CREATE TABLE organization_role_permission_links (
    role_id        UUID NOT NULL REFERENCES organization_roles (id) ON DELETE CASCADE,
    permission_code VARCHAR(255) NOT NULL REFERENCES organization_role_permissions (code) ON DELETE CASCADE,

    PRIMARY KEY (role_id, permission_code)
);

-- +migrate Down
DROP TABLE IF EXISTS organization_role_permission_links CASCADE;
DROP TABLE IF EXISTS organization_role_permissions CASCADE;
DROP TABLE IF EXISTS organization_member_roles CASCADE;
DROP TABLE IF EXISTS organization_roles CASCADE;

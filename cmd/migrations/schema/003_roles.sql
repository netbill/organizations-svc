-- +migrate Up
CREATE TABLE roles (
    id               UUID    PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    organization_id  UUID    NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    head             BOOLEAN NOT NULL DEFAULT false,
    rank             INT     NOT NULL DEFAULT 0 CHECK (rank >= 0),
    name             TEXT    NOT NULL,
    description      TEXT    NOT NULL,
    color            TEXT    NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (organization_id, name)
);

CREATE UNIQUE INDEX roles_one_head_per_organization
    ON roles (organization_id)
    WHERE head = true;

CREATE TABLE member_roles (
    member_id UUID NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    role_id   UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,

    PRIMARY KEY (member_id, role_id)
);

-- permissions dictionary
CREATE TABLE role_permissions (
    id          UUID          PRIMARY KEY,
    code        VARCHAR(255)  UNIQUE NOT NULL,
    description VARCHAR(1024) NOT NULL
);

INSERT INTO role_permissions (id, code, description) VALUES
    (uuid_generate_v4(), 'organization.manage', 'manage organization settings'),
    (uuid_generate_v4(), 'invites.manage', 'manage organization invites'),
    (uuid_generate_v4(), 'members.manage', 'manage organization members'),
    (uuid_generate_v4(), 'roles.manage', 'manage organization roles');

-- role ↔ permission links
CREATE TABLE role_permission_links (
    role_id       UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES role_permissions(id) ON DELETE CASCADE,

    PRIMARY KEY (role_id, permission_id)
);

-- 1) if role.head=true -> add all permissions to role
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION ensure_head_role_permissions()
RETURNS trigger AS $$
BEGIN
    IF NEW.head = true THEN
        INSERT INTO role_permission_links (role_id, permission_id)
        SELECT NEW.id, p.id
        FROM role_permissions p
        ON CONFLICT DO NOTHING;
    END IF;

    RETURN NEW;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_roles_ensure_head_perms_ins ON roles;
CREATE TRIGGER trg_roles_ensure_head_perms_ins
AFTER INSERT ON roles
FOR EACH ROW
EXECUTE FUNCTION ensure_head_role_permissions();

DROP TRIGGER IF EXISTS trg_roles_ensure_head_perms_upd ON roles;
CREATE TRIGGER trg_roles_ensure_head_perms_upd
AFTER UPDATE OF head ON roles
FOR EACH ROW
EXECUTE FUNCTION ensure_head_role_permissions();

CREATE OR REPLACE FUNCTION grant_new_permission_to_head_roles()
RETURNS trigger AS $$
BEGIN
    INSERT INTO role_permission_links (role_id, permission_id)
    SELECT r.id, NEW.id
    FROM roles r
    WHERE r.head = true
    ON CONFLICT DO NOTHING;

RETURN NEW;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_role_permissions_grant_to_head_roles ON role_permissions;
CREATE TRIGGER trg_role_permissions_grant_to_head_roles
AFTER INSERT ON role_permissions
FOR EACH ROW
EXECUTE FUNCTION grant_new_permission_to_head_roles();
-- +migrate StatementEnd

-- 3) ban delete permissions from head-roles
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION prevent_delete_head_role_permissions()
RETURNS trigger AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM roles r
        WHERE r.id = OLD.role_id
        AND r.head = true
    ) THEN
        RAISE EXCEPTION 'cannot delete permissions from head role %', OLD.role_id
            USING ERRCODE = '23514';
    END IF;

    RETURN OLD;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_role_permission_links_prevent_delete_head ON role_permission_links;
CREATE TRIGGER trg_role_permission_links_prevent_delete_head
BEFORE DELETE ON role_permission_links
FOR EACH ROW
EXECUTE FUNCTION prevent_delete_head_role_permissions();
-- +migrate StatementEnd

-- 4) ban change of organization_id for roles
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION prevent_role_organization_change()
RETURNS trigger AS $$
BEGIN
    IF NEW.organization_id <> OLD.organization_id THEN
        RAISE EXCEPTION 'cannot change organization_id for role %', OLD.id
            USING ERRCODE = '23514';
    END IF;

    RETURN NEW;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_roles_prevent_organization_change ON roles;
CREATE TRIGGER trg_roles_prevent_organization_change
BEFORE UPDATE OF organization_id ON roles
FOR EACH ROW
EXECUTE FUNCTION prevent_role_organization_change();
-- +migrate StatementEnd

-- 5) ban delete head-roles
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION prevent_delete_head_role()
RETURNS trigger AS $$
BEGIN
    IF OLD.head = true THEN
        RAISE EXCEPTION 'cannot delete head role %', OLD.id
            USING ERRCODE = '23514';
    END IF;

    RETURN OLD;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_roles_prevent_delete_head ON roles;
CREATE TRIGGER trg_roles_prevent_delete_head
BEFORE DELETE ON roles
FOR EACH ROW
EXECUTE FUNCTION prevent_delete_head_role();
-- +migrate StatementEnd

-- 6) ban remove head-role from members
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION prevent_remove_head_role_from_member()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM roles r
        WHERE r.id = OLD.role_id
        AND r.head = true
    ) THEN
        RAISE EXCEPTION 'cannot remove head role % from member %',
        OLD.role_id, OLD.member_id
        USING ERRCODE = '23514';
    END IF;

    RETURN OLD;
END;
$$;
-- +migrate StatementEnd

DROP TRIGGER IF EXISTS trg_member_roles_prevent_delete_head_role ON member_roles;
CREATE TRIGGER trg_member_roles_prevent_delete_head_role
BEFORE DELETE ON member_roles
FOR EACH ROW
EXECUTE FUNCTION prevent_remove_head_role_from_member();

-- +migrate Down
DROP TRIGGER IF EXISTS trg_member_roles_prevent_delete_head_role ON member_roles;
DROP FUNCTION IF EXISTS prevent_remove_head_role_from_member();

DROP TRIGGER IF EXISTS trg_roles_prevent_delete_head ON roles;
DROP FUNCTION IF EXISTS prevent_delete_head_role();

DROP TRIGGER IF EXISTS trg_roles_prevent_organization_change ON roles;
DROP FUNCTION IF EXISTS prevent_role_organization_change();

DROP TRIGGER IF EXISTS trg_role_permission_links_prevent_delete_head ON role_permission_links;
DROP FUNCTION IF EXISTS prevent_delete_head_role_permissions();

DROP TRIGGER IF EXISTS trg_role_permissions_grant_to_head_roles ON role_permissions;
DROP FUNCTION IF EXISTS grant_new_permission_to_head_roles();

DROP TRIGGER IF EXISTS trg_roles_ensure_head_perms_upd ON roles;
DROP TRIGGER IF EXISTS trg_roles_ensure_head_perms_ins ON roles;
DROP FUNCTION IF EXISTS ensure_head_role_permissions();

DROP TABLE IF EXISTS role_permission_links CASCADE;
DROP TABLE IF EXISTS role_permissions CASCADE;
DROP TABLE IF EXISTS member_roles CASCADE;
DROP TABLE IF EXISTS roles CASCADE;

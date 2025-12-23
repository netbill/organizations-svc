-- +migrate Up


-- 1) if role.head=true -> add all permissions to role_permissions

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION ensure_head_role_permissions()
RETURNS trigger AS $$
BEGIN
    IF NEW.head = true THEN
        INSERT INTO role_permissions (role_id, permission_id)
        SELECT NEW.id, p.id
        FROM permissions p
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
    INSERT INTO role_permissions (role_id, permission_id)
    SELECT r.id, NEW.id
    FROM roles r
    WHERE r.head = true
    ON CONFLICT DO NOTHING;

    RETURN NEW;
END
$$ LANGUAGE plpgsql;


DROP TRIGGER IF EXISTS trg_permissions_grant_to_head_roles ON permissions;
CREATE TRIGGER trg_permissions_grant_to_head_roles
AFTER INSERT ON permissions
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

DROP TRIGGER IF EXISTS trg_role_permissions_prevent_delete_head ON role_permissions;
CREATE TRIGGER trg_role_permissions_prevent_delete_head
BEFORE DELETE ON role_permissions
FOR EACH ROW
EXECUTE FUNCTION prevent_delete_head_role_permissions();
-- +migrate StatementEnd

-- 4) ban change of agglomeration_id for roles

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION prevent_role_agglomeration_change()
RETURNS trigger AS $$
BEGIN
    IF NEW.agglomeration_id <> OLD.agglomeration_id THEN
        RAISE EXCEPTION 'cannot change agglomeration_id for role %', OLD.id
            USING ERRCODE = '23514';
    END IF;

    RETURN NEW;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_roles_prevent_agglomeration_change ON roles;
CREATE TRIGGER trg_roles_prevent_agglomeration_change
BEFORE UPDATE OF agglomeration_id ON roles
FOR EACH ROW
EXECUTE FUNCTION prevent_role_agglomeration_change();
-- +migrate StatementEnd

-- 5) ban delete head-roles (roles)

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

-- 6) ban remove head-roles from members

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION prevent_remove_head_role_from_member() RETURNS trigger
    LANGUAGE plpgsql
AS $$ BEGIN
        IF EXISTS (
            SELECT 1
            FROM roles r
            WHERE r.id = OLD.role_id
            AND r.head = true
        ) THEN
            RAISE EXCEPTION 'cannot remove head role % from member %', OLD.role_id, OLD.member_id
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

DROP TRIGGER IF EXISTS trg_roles_prevent_agglomeration_change ON roles;
DROP FUNCTION IF EXISTS prevent_role_agglomeration_change();

DROP TRIGGER IF EXISTS trg_role_permissions_prevent_delete_head ON role_permissions;
DROP FUNCTION IF EXISTS prevent_delete_head_role_permissions();

DROP TRIGGER IF EXISTS trg_permissions_grant_to_head_roles ON permissions;
DROP FUNCTION IF EXISTS grant_new_permission_to_head_roles();

DROP TRIGGER IF EXISTS trg_roles_ensure_head_perms_upd ON roles;
DROP TRIGGER IF EXISTS trg_roles_ensure_head_perms_ins ON roles;
DROP FUNCTION IF EXISTS ensure_head_role_permissions();

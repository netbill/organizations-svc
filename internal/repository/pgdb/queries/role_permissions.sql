-- name: CreateRolePermission :exec
INSERT INTO role_permissions (
    role_id,
    permission_id
) VALUES (
    $1, $2
);

-- name: DeleteRolePermission :exec
DELETE FROM role_permissions
WHERE role_id = $1
AND permission_id = $2;

-- name: UserHasPermission :one
SELECT EXISTS (
    SELECT 1
    FROM members m
             JOIN member_roles mr ON mr.member_id = m.id
             JOIN role_permissions rp ON rp.role_id = mr.role_id
             JOIN permissions p ON p.id = rp.permission_id
    WHERE m.account_id = $1
      AND m.agglomeration_id = $2
      AND p.code = $3
);

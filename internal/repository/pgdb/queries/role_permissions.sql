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

-- name: CheckMemberHavePermissionByCode :one
SELECT EXISTS (
    SELECT 1
    FROM members m
    JOIN member_roles mr ON mr.member_id = m.id
    JOIN role_permissions rp ON rp.role_id = mr.role_id
    JOIN permissions p ON p.id = rp.permission_id
    WHERE m.id = sqlc.arg('member_id')::uuid
        AND m.agglomeration_id = sqlc.arg('agglomeration_id')::uuid
        AND p.code = sqlc.arg('code')::text
);


-- name: CheckMemberHavePermissionByID :one
SELECT EXISTS (
    SELECT 1
    FROM members m
    JOIN member_roles mr ON mr.member_id = m.id
    JOIN role_permissions rp ON rp.role_id = mr.role_id
    WHERE m.id = sqlc.arg('member_id')::uuid
        AND m.agglomeration_id = sqlc.arg('agglomeration_id')::uuid
        AND rp.permission_id = sqlc.arg('permission_id')::uuid
);

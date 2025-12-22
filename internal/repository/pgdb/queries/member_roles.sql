-- name: CreateMemberRole :exec
INSERT INTO member_roles (
    member_id,
    role_id
) VALUES (
    sqlc.arg('member_id')::uuid,
    sqlc.arg('role_id')::uuid
)
ON CONFLICT DO NOTHING;

-- name: DeleteMemberRole :exec
DELETE FROM member_roles
WHERE member_id = sqlc.arg('member_id')::uuid
  AND role_id = sqlc.arg('role_id')::uuid;

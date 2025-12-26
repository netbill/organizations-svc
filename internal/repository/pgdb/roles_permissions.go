package pgdb

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/pgx"
)

type RolePermissionsQ struct {
	db pgx.DBTX
}

func NewRolePermissionsQ(db pgx.DBTX) RolePermissionsQ {
	return RolePermissionsQ{db: db}
}

func (q RolePermissionsQ) New() RolePermissionsQ { return NewRolePermissionsQ(q.db) }

func (q RolePermissionsQ) Insert(
	ctx context.Context,
	roleID uuid.UUID,
	permissionID uuid.UUID,
) error {
	const sqlq = `
		INSERT INTO role_permissions (role_id, permission_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`

	_, err := q.db.ExecContext(ctx, sqlq, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("inserting role_permission: %w", err)
	}
	return nil
}

func (q RolePermissionsQ) Delete(
	ctx context.Context,
	roleID uuid.UUID,
	permissionID uuid.UUID,
) error {
	const sqlq = `
		DELETE FROM role_permissions
		WHERE role_id = $1 AND permission_id = $2
	`

	_, err := q.db.ExecContext(ctx, sqlq, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("deleting role_permission: %w", err)
	}
	return nil
}

func (q RolePermissionsQ) CheckMemberHavePermissionByCode(
	ctx context.Context,
	memberID uuid.UUID,
	code string,
) (bool, error) {
	const sqlq = `
		SELECT EXISTS (
			SELECT 1
			FROM members m
			JOIN member_roles mr ON mr.member_id = m.id
			JOIN role_permissions rp ON rp.role_id = mr.role_id
			JOIN permissions p ON p.id = rp.permission_id
			WHERE m.id = $1
				AND p.code = $2
		)
	`

	var ok bool
	if err := q.db.QueryRowContext(ctx, sqlq, memberID, code).Scan(&ok); err != nil {
		return false, fmt.Errorf("scanning permission exists (member+code): %w", err)
	}
	return ok, nil
}

func (q RolePermissionsQ) CheckMemberHavePermissionByID(
	ctx context.Context,
	memberID uuid.UUID,
	permissionID uuid.UUID,
) (bool, error) {
	const sqlq = `
		SELECT EXISTS (
			SELECT 1
			FROM members m
			JOIN member_roles mr ON mr.member_id = m.id
			JOIN role_permissions rp ON rp.role_id = mr.role_id
			WHERE m.id = $1
				AND rp.permission_id = $2
		)
	`

	var ok bool
	if err := q.db.QueryRowContext(ctx, sqlq, memberID, permissionID).Scan(&ok); err != nil {
		return false, fmt.Errorf("scanning permission exists (member+id): %w", err)
	}
	return ok, nil
}

func (q RolePermissionsQ) CheckAccountHavePermissionByCode(
	ctx context.Context,
	accountID uuid.UUID,
	agglomerationID uuid.UUID,
	code string,
) (bool, error) {
	const sqlq = `
		SELECT EXISTS (
			SELECT 1
			FROM members m
			JOIN member_roles mr ON mr.member_id = m.id
			JOIN role_permissions rp ON rp.role_id = mr.role_id
			JOIN permissions p ON p.id = rp.permission_id
			WHERE m.account_id = $1
				AND m.agglomeration_id = $2
				AND p.code = $3
		)
	`

	var ok bool
	if err := q.db.QueryRowContext(ctx, sqlq, accountID, agglomerationID, code).Scan(&ok); err != nil {
		return false, fmt.Errorf("scanning permission exists (account+code): %w", err)
	}
	return ok, nil
}

func (q RolePermissionsQ) CheckAccountHavePermissionByID(
	ctx context.Context,
	accountID uuid.UUID,
	agglomerationID uuid.UUID,
	permissionID uuid.UUID,
) (bool, error) {
	const sqlq = `
		SELECT EXISTS (
			SELECT 1
			FROM members m
			JOIN member_roles mr ON mr.member_id = m.id
			JOIN role_permissions rp ON rp.role_id = mr.role_id
			WHERE m.account_id = $1
				AND m.agglomeration_id = $2
				AND rp.permission_id = $3
		)
	`

	var ok bool
	if err := q.db.QueryRowContext(ctx, sqlq, accountID, agglomerationID, permissionID).Scan(&ok); err != nil {
		return false, fmt.Errorf("scanning permission exists (account+id): %w", err)
	}
	return ok, nil
}

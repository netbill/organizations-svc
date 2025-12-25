package pgdbsq

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
)

type PermissionCheckQ struct {
	db pgdb.DBTX
}

func NewPermissionCheckQ(db pgdb.DBTX) PermissionCheckQ {
	return PermissionCheckQ{db: db}
}

func (q PermissionCheckQ) New() PermissionCheckQ { return NewPermissionCheckQ(q.db) }

func (q PermissionCheckQ) CheckMemberHavePermissionByCode(
	ctx context.Context,
	memberID uuid.UUID,
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
			WHERE m.id = $1
				AND m.agglomeration_id = $2
				AND p.code = $3
		)
	`

	var ok bool
	if err := q.db.QueryRowContext(ctx, sqlq, memberID, agglomerationID, code).Scan(&ok); err != nil {
		return false, fmt.Errorf("scanning permission exists (member+code): %w", err)
	}
	return ok, nil
}

func (q PermissionCheckQ) CheckMemberHavePermissionByID(
	ctx context.Context,
	memberID uuid.UUID,
	agglomerationID uuid.UUID,
	permissionID uuid.UUID,
) (bool, error) {
	const sqlq = `
		SELECT EXISTS (
			SELECT 1
			FROM members m
			JOIN member_roles mr ON mr.member_id = m.id
			JOIN role_permissions rp ON rp.role_id = mr.role_id
			WHERE m.id = $1
				AND m.agglomeration_id = $2
				AND rp.permission_id = $3
		)
	`

	var ok bool
	if err := q.db.QueryRowContext(ctx, sqlq, memberID, agglomerationID, permissionID).Scan(&ok); err != nil {
		return false, fmt.Errorf("scanning permission exists (member+id): %w", err)
	}
	return ok, nil
}

func (q PermissionCheckQ) CheckAccountHavePermissionByCode(
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

func (q PermissionCheckQ) CheckAccountHavePermissionByID(
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

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
)

func (s Service) CreatePermissionForRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	return s.sql(ctx).CreateRolePermission(ctx, pgdb.CreateRolePermissionParams{
		RoleID:       roleID,
		PermissionID: permissionID,
	})
}

func (s Service) DeletePermissionForRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	return s.sql(ctx).DeleteRolePermission(ctx, pgdb.DeleteRolePermissionParams{
		RoleID:       roleID,
		PermissionID: permissionID,
	})
}

func (s Service) CheckMemberHavePermissionByID(ctx context.Context, memberID, permissionID uuid.UUID) (bool, error) {
	has, err := s.sql(ctx).CheckMemberHavePermissionByID(ctx, pgdb.CheckMemberHavePermissionByIDParams{
		MemberID:     memberID,
		PermissionID: permissionID,
	})
	if err != nil {
		return false, err
	}

	return has, nil
}

func (s Service) CheckMemberHavePermissionByCode(ctx context.Context, memberID uuid.UUID, permissionKey string) (bool, error) {
	has, err := s.sql(ctx).CheckMemberHavePermissionByCode(ctx, pgdb.CheckMemberHavePermissionByCodeParams{
		MemberID: memberID,
		Code:     permissionKey,
	})
	if err != nil {
		return false, err
	}

	return has, nil
}

func (s Service) CheckAccountHavePermissionByCode(ctx context.Context, accountID, agglomerationID uuid.UUID, permissionKey string) (bool, error) {
	has, err := s.sql(ctx).CheckAccountHavePermissionByCode(ctx, pgdb.CheckAccountHavePermissionByCodeParams{
		AccountID:       accountID,
		AgglomerationID: agglomerationID,
		Code:            permissionKey,
	})
	if err != nil {
		return false, err
	}

	return has, nil
}

func (s Service) CheckAccountHavePermissionByID(ctx context.Context, accountID, agglomerationID, permissionID uuid.UUID) (bool, error) {
	has, err := s.sql(ctx).CheckAccountHavePermissionByID(ctx, pgdb.CheckAccountHavePermissionByIDParams{
		AccountID:       accountID,
		AgglomerationID: agglomerationID,
		PermissionID:    permissionID,
	})
	if err != nil {
		return false, err
	}

	return has, nil
}

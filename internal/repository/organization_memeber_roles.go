package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/repository/pgdb"
)

func (r Repository) GetMemberRoles(ctx context.Context, memberID uuid.UUID) ([]models.Role, error) {
	memberRoles, err := r.orgRolesQ(ctx).
		FilterByMemberID(memberID).
		OrderByRoleRank(true).
		Select(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles for member ID %s, cause: %w", memberID, err)
	}

	result := make([]models.Role, len(memberRoles))
	for i, mr := range memberRoles {
		result[i] = Role(mr)
	}

	return result, nil
}

func (r Repository) RemoveMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error {
	err := r.orgMemberRolesQ(ctx).
		FilterByMemberID(memberID).
		FilterByRoleID(roleID).
		Delete(ctx)

	if err != nil {
		return fmt.Errorf("failed to remove role ID %s from member ID %s, cause: %w", roleID, memberID, err)
	}

	return nil
}

func (r Repository) AddMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error {
	_, err := r.orgMemberRolesQ(ctx).
		Insert(ctx, pgdb.OrganizationMemberRole{
			MemberID: memberID,
			RoleID:   roleID,
		})

	if err != nil {
		return fmt.Errorf("failed to add role ID %s to member ID %s, cause: %w", roleID, memberID, err)
	}

	return nil
}

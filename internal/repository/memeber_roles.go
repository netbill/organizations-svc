package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/repository/pgdb"
)

func (s Service) GetMemberRoles(ctx context.Context, memberID uuid.UUID) ([]models.Role, error) {
	memberRoles, err := s.rolesQ(ctx).
		FilterByMemberID(memberID).
		OrderByRoleRank(true).
		Select(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]models.Role, len(memberRoles))
	for i, mr := range memberRoles {
		result[i] = Role(mr)
	}

	return result, nil
}

func (s Service) RemoveMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error {
	return s.memberRolesQ(ctx).
		FilterByMemberID(memberID).
		FilterByRoleID(roleID).
		Delete(ctx)
}

func (s Service) AddMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error {
	_, err := s.memberRolesQ(ctx).
		Insert(ctx, pgdb.MemberRole{
			MemberID: memberID,
			RoleID:   roleID,
		})

	return err
}

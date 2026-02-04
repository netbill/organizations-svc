package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/netbill/organizations-svc/internal/core/models"
)

type OrganizationMemberRolesRow struct {
	MemberID  uuid.UUID `db:"member_id"`
	RoleID    uuid.UUID `db:"role_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (r OrganizationMemberRolesRow) IsNil() bool {
	return r.MemberID == uuid.Nil && r.RoleID == uuid.Nil
}

func (r OrganizationMemberRolesRow) ToModel() models.OrgMemberRolesLink {
	return models.OrgMemberRolesLink{
		MemberID:  r.MemberID,
		RoleID:    r.RoleID,
		CreatedAt: r.CreatedAt,
	}
}

type OrgMemberRolesQ interface {
	New() OrgMemberRolesQ

	Insert(ctx context.Context, input OrganizationMemberRolesRow) (OrganizationMemberRolesRow, error)

	Delete(ctx context.Context) error
	FilterByMemberID(memberID uuid.UUID) OrgMemberRolesQ
	FilterByRoleID(roleID uuid.UUID) OrgMemberRolesQ

	Select(ctx context.Context) ([]OrganizationMemberRolesRow, error)
	Get(ctx context.Context) (OrganizationMemberRolesRow, error)
}

func (r *Repository) GetMemberRoles(ctx context.Context, memberID uuid.UUID) ([]models.Role, error) {
	memberRoles, err := r.orgRolesQ().
		FilterByMemberID(memberID).
		OrderByRoleRank(true).
		Select(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get member roles, cause: %w", err)
	}

	result := make([]models.Role, len(memberRoles))
	for i, mr := range memberRoles {
		result[i] = mr.ToModel()
	}

	return result, nil
}

func (r *Repository) RemoveMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error {
	err := r.orgMemberRolesQ().
		FilterByMemberID(memberID).
		FilterByRoleID(roleID).
		Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to remove member role, cause: %w", err)
	}

	return nil
}

func (r *Repository) AddMemberRole(
	ctx context.Context,
	memberID, roleID uuid.UUID,
) (models.OrgMemberRolesLink, error) {
	link, err := r.orgMemberRolesQ().
		Insert(ctx, OrganizationMemberRolesRow{
			MemberID: memberID,
			RoleID:   roleID,
		})

	if err != nil {
		return models.OrgMemberRolesLink{}, fmt.Errorf("failed to add member role, cause: %w", err)
	}

	return link.ToModel(), nil
}

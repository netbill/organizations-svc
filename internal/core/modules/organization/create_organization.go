package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type CreateParams struct {
	Name string
}

func (m *Module) CreateOrganization(
	ctx context.Context,
	accountID uuid.UUID,
	params CreateParams,
) (org models.Organization, err error) {
	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		org, err = m.repo.CreateOrganization(ctx, params)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrganizationCreated(ctx, org)
		if err != nil {
			return err
		}

		role, err := m.createRoleHead(ctx, org.ID)
		if err != nil {
			return err
		}

		if _, err = m.createMemberHead(ctx, accountID, org.ID, role.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Organization{}, err
	}

	return org, err
}

func (m *Module) createRoleHead(
	ctx context.Context,
	organizationID uuid.UUID,
) (role models.Role, err error) {
	role, err = m.repo.CreateHeadRole(ctx, organizationID)
	if err != nil {
		return models.Role{}, err
	}

	err = m.messenger.WriteOrgRoleCreated(ctx, role)
	if err != nil {
		return models.Role{}, err
	}

	per, err := m.repo.GetRolePermissions(ctx, role.ID)
	if err != nil {
		return models.Role{}, err
	}

	if err = m.messenger.WriteOrgRolePermissionsUpdated(ctx, role, per); err != nil {
		return models.Role{}, err
	}

	return role, nil
}

func (m *Module) createMemberHead(
	ctx context.Context,
	accountID uuid.UUID,
	organizationID uuid.UUID,
	roleID uuid.UUID,
) (member models.Member, err error) {
	member, err = m.repo.CreateMember(ctx, accountID, organizationID)
	if err != nil {
		return models.Member{}, err
	}

	err = m.repo.AddMemberRole(ctx, member.ID, roleID)
	if err != nil {
		return models.Member{}, err
	}

	err = m.messenger.WriteOrgMemberCreated(ctx, member)
	if err != nil {
		return models.Member{}, err
	}

	return member, nil
}

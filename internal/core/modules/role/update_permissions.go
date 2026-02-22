package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
	"github.com/netbill/orgperm"
)

type SetPermissions map[uuid.UUID]bool

func (s SetPermissions) Validate() error {
	if len(s) != len(orgperm.GetAllPermissions()) {
		return fmt.Errorf("permission count mismatch")
	}

	return nil
}

func (m *Module) UpdatePermissions(
	ctx context.Context,
	initiator domain.AccountActor,
	roleID uuid.UUID,
	params SetPermissions,
) (role domain.Role, links domain.OrgRolePermissionsWithDetailsForRole, err error) {
	if err = params.Validate(); err != nil {
		return domain.Role{}, domain.OrgRolePermissionsWithDetailsForRole{}, err
	}

	role, err = m.repo.GetRole(ctx, roleID)
	if err != nil {
		return domain.Role{}, nil, err
	}

	member, err := m.getInitiator(ctx, initiator, role.OrganizationID)
	if err != nil {
		return domain.Role{}, nil, err
	}

	err = m.checkPermissionsToManageRole(ctx, initiator, member.OrganizationID, role.Rank)
	if err != nil {
		return domain.Role{}, nil, err
	}

	err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		links, err = m.repo.SetRolePermissions(ctx, roleID, params)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrgRolePermissionsUpdated(ctx, role, params)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return domain.Role{}, nil, err
	}

	return role, links, nil
}

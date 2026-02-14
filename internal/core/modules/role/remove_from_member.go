package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) RemoveFromMember(
	ctx context.Context,
	initiator models.AccountActor,
	memberID, roleID uuid.UUID,
) error {
	member, err := m.getMember(ctx, memberID)
	if err != nil {
		return err
	}

	role, err := m.GetByID(ctx, roleID)
	if err != nil {
		return err
	}

	if role.OrganizationID != member.OrganizationID {
		return errx.ErrorRoleNotFound.Raise(
			fmt.Errorf("role with id %s is not available in organization %s", role.ID, role.OrganizationID),
		)
	}

	err = m.checkPermissionsToManageRole(ctx, initiator, role.OrganizationID, role.Rank)
	if err != nil {
		return err
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		if err = m.repo.RemoveMemberRole(ctx, memberID, roleID); err != nil {
			return err
		}

		if err = m.messenger.WriteOrgMemberRoleRemove(ctx, memberID, roleID); err != nil {
			return err
		}

		return nil
	})
}

package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) AddForMember(
	ctx context.Context,
	initiator models.AccountActor,
	memberID, roleID uuid.UUID,
) error {
	role, err := m.checkPermissionsToManageRole(ctx, initiator, roleID)
	if err != nil {
		return err
	}

	member, err := m.repo.GetMember(ctx, memberID)
	if err != nil {
		return err
	}

	if member.OrganizationID != role.OrganizationID {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s does not belong to organization %s", memberID, role.OrganizationID),
		)
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		if err = m.repo.LockMemberRolesLinksRevision(ctx, memberID); err != nil {
			return err
		}

		res, err := m.repo.AddMemberRole(ctx, memberID, roleID)
		if err != nil {
			return err
		}

		revision, err := m.repo.BumpMemberRolesLinksRevision(ctx, memberID)
		if err != nil {
			return err
		}

		if err = m.messenger.WriteOrgMemberRolesUpdated(ctx, memberID, res, revision); err != nil {
			return err
		}

		return nil
	})
}

func (m *Module) RemoveFromMember(
	ctx context.Context,
	initiator models.AccountActor,
	memberID, roleID uuid.UUID,
) error {
	role, err := m.checkPermissionsToManageRole(ctx, initiator, roleID)
	if err != nil {
		return err
	}

	member, err := m.repo.GetMember(ctx, memberID)
	if err != nil {
		return err
	}

	if member.OrganizationID != role.OrganizationID {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s does not belong to organization %s", memberID, role.OrganizationID),
		)
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		err = m.repo.LockMemberRolesLinksRevision(ctx, memberID)
		if err != nil {
			return err
		}

		roles, err := m.repo.RemoveMemberRole(ctx, memberID, roleID)
		if err != nil {
			return err
		}

		revision, err := m.repo.BumpMemberRolesLinksRevision(ctx, memberID)
		if err != nil {
			return err
		}

		if err = m.messenger.WriteOrgMemberRolesUpdated(ctx, memberID, roles, revision); err != nil {
			return err
		}

		return nil
	})
}

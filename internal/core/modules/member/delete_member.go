package member

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) DeleteMember(
	ctx context.Context,
	initiator models.InitiatorData,
	memberID uuid.UUID,
) error {
	member, err := m.GetMemberByID(ctx, memberID)
	if err != nil {
		return err
	}
	if member.Head {
		return errx.ErrorCannotDeleteOrganizationHeadMember.Raise(
			fmt.Errorf("cannot delete organization head member %s", member.ID),
		)
	}

	err = m.checkPermissionToInteractWithMember(ctx, initiator, member.OrganizationID, memberID)
	if err != nil {
		return err
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		err = m.repo.DeleteMember(ctx, memberID)
		if err != nil {
			return fmt.Errorf("failed to delete member %s: %w", memberID, err)
		}

		if err = m.messenger.WriteOrgMemberDeleted(ctx, member.ID); err != nil {
			return fmt.Errorf("failed to send member deleted message for member %s: %w", memberID, err)
		}

		return nil
	})
}

package member

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) Delete(
	ctx context.Context,
	actor models.AccountActor,
	memberID uuid.UUID,
) error {
	if actor == memberID {
		return errx.ErrorCannotDeleteSelf.Raise(
			fmt.Errorf("member %s cannot delete itself", actor),
		)
	}

	member, err := m.GetByID(ctx, actor)
	if err != nil {
		return err
	}
	if member.Head {
		return errx.ErrorCannotDeleteOrganizationHeadMember.Raise(
			fmt.Errorf("cannot delete organization head member %s", member.ID),
		)
	}

	initiator, err := m.repo.GetMemberByAccountAndOrganization(ctx, actor, member.OrganizationID)
	if err != nil {
		return err
	}
	if !initiator.Head {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("account has no rights to delete member %s", memberID),
		)
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		err = m.repo.DeleteMember(ctx, memberID)
		if err != nil {
			return fmt.Errorf("failed to delete member %s: %w", memberID, err)
		}

		if err = m.messenger.WriteOrgMemberDeleted(ctx, memberID); err != nil {
			return fmt.Errorf("failed to send member deleted message for member %s: %w", memberID, err)
		}

		return nil
	})
}

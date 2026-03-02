package member

import (
	"context"
	"errors"
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
	member, err := m.GetByID(ctx, memberID)
	if errors.Is(err, errx.ErrorMemberNotExists) {
		buried, err := m.repo.MemberIsBuried(ctx, memberID)
		if err != nil {
			return err
		}
		if buried {
			return errx.ErrorMemberDeleted.Raise(
				fmt.Errorf("member with id %s is already deleted", memberID),
			)
		}
	}
	if err != nil {
		return err
	}
	if member.Head {
		return errx.ErrorCannotDeleteOrganizationHeadMember.Raise(
			fmt.Errorf("cannot delete organization head member %s", member.ID),
		)
	}

	if member.AccountID == actor {
		return errx.ErrorCannotDeleteSelf.Raise(
			fmt.Errorf("account cannot delete itself as member %s", member.ID),
		)
	}

	organization, err := m.repo.GetOrganizationByID(ctx, member.OrganizationID)
	if err != nil {
		return err
	}
	if organization.Status == models.OrganizationStatusSuspended {
		return errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization %s is suspended", member.OrganizationID),
		)
	}

	initiator, err := m.getInitiator(ctx, actor, member.OrganizationID)
	if err != nil {
		return err
	}
	if !initiator.Head {
		return errx.ErrorNotOrganizationHead.Raise(
			fmt.Errorf("account has no rights to delete member %s", memberID),
		)
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		err = m.repo.BuryMember(ctx, memberID)
		if err != nil {
			return fmt.Errorf("failed to bury member %s: %w", memberID, err)
		}

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

func (m *Module) DeleteSelf(ctx context.Context, actor models.AccountActor, orgID uuid.UUID) error {
	_, err := m.repo.GetOrganizationByID(ctx, orgID)
	if err != nil {
		return err
	}

	member, err := m.getInitiator(ctx, actor, orgID)
	if err != nil {
		return err
	}
	if member.Head {
		return errx.ErrorCannotDeleteOrganizationHeadMember.Raise(
			fmt.Errorf("cannot delete organization head member %s", member.ID),
		)
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		if err = m.repo.BuryMember(ctx, member.ID); err != nil {
			return fmt.Errorf("failed to bury member %s: %w", member.ID, err)
		}

		if err = m.repo.DeleteMember(ctx, member.ID); err != nil {
			return fmt.Errorf("failed to delete member %s: %w", member.ID, err)
		}

		if err = m.messenger.WriteOrgMemberDeleted(ctx, member.ID); err != nil {
			return fmt.Errorf("failed to send member deleted message for member %s: %w", member.ID, err)
		}

		return nil
	})
}

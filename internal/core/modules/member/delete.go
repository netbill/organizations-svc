package member

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/orgperm"
)

func (m *Module) Delete(
	ctx context.Context,
	initiator domain.AccountActor,
	memberID uuid.UUID,
) error {
	member, err := m.GetByID(ctx, memberID)
	if err != nil {
		return err
	}
	if member.Head {
		return errx.ErrorCannotDeleteOrganizationHeadMember.Raise(
			fmt.Errorf("cannot delete organization head member %s", member.ID),
		)
	}

	err = m.checkAbilityToDeleteMember(ctx, initiator, member.OrganizationID, memberID)
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

func (m *Module) checkAbilityToDeleteMember(
	ctx context.Context,
	initiator domain.AccountActor,
	organizationID uuid.UUID,
	memberID uuid.UUID,
) error {
	member, err := m.GetInitiator(ctx, initiator, organizationID)
	if err != nil {
		return err
	}

	if !member.Head {
		hasPermission, err := m.repo.CheckMemberHavePermission(ctx, member.AccountID, orgperm.MembersDeleteID)
		if err != nil {
			return err
		}
		if !hasPermission {
			return errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("initiator member %s has no delete members permission", member.ID),
			)
		}

		firstMaxRole, err := m.repo.GetMemberMaxRole(ctx, member.ID)
		if err != nil {
			return fmt.Errorf("failed to get max role for member %s: %w", member.AccountID, err)
		}

		secMaxRole, err := m.repo.GetMemberMaxRole(ctx, memberID)
		if err != nil {
			return fmt.Errorf("failed to get max role for member %s: %w", memberID, err)
		}

		if firstMaxRole.Rank < secMaxRole.Rank {
			return errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf(
					"member %s with rank %d cannot manage member %s with rank %d",
					member.AccountID,
					firstMaxRole.Rank,
					memberID,
					secMaxRole.Rank,
				),
			)
		}
	}

	return nil
}

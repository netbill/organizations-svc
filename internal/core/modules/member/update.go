package member

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type UpdateParams struct {
	Position *string
	Label    *string
}

func (m *Module) Update(
	ctx context.Context,
	initiator models.AccountActor,
	memberID uuid.UUID,
	params UpdateParams,
) (models.Member, error) {
	member, err := m.GetByID(ctx, memberID)
	if err != nil {
		return models.Member{}, err
	}
	if member.Head {
		return models.Member{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("cannot update organization head member %s", member.ID),
		)
	}

	err = m.checkAbilityToUpdateMember(ctx, initiator, member.OrganizationID, memberID)
	if err != nil {
		return models.Member{}, err
	}

	err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		member, err = m.repo.UpdateMember(ctx, memberID, params)
		if err != nil {
			return fmt.Errorf("failed to update member %s: %w", memberID, err)
		}

		if err = m.messenger.WriteOrgMemberUpdated(ctx, member); err != nil {
			return fmt.Errorf("failed to send member updated message for member %s: %w", memberID, err)
		}
		return nil
	})

	return member, err
}

package member

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/domain/errx"
	"github.com/netbill/organizations-svc/internal/domain/models"
)

type UpdateParams struct {
	Position *string
	Label    *string
}

func (s Service) UpdateMember(
	ctx context.Context,
	accountID, memberID uuid.UUID,
	params UpdateParams,
) (models.Member, error) {
	member, err := s.GetMemberByID(ctx, memberID)
	if err != nil {
		return models.Member{}, err
	}

	initiator, err := s.GetInitiatorMember(ctx, accountID, memberID)
	if err != nil {
		return models.Member{}, err
	}

	if err = s.CheckAccessToManageOtherMember(ctx, initiator.ID, member.ID); err != nil {
		return models.Member{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member %s has no permission to manage members: %w", initiator.ID, err),
		)
	}

	err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		member, err = s.repo.UpdateMember(ctx, memberID, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update member %s: %w", memberID, err),
			)
		}

		if err = s.messenger.WriteMemberUpdated(ctx, member); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send member updated message for member %s: %w", memberID, err),
			)
		}
		return nil
	})

	return member, err
}

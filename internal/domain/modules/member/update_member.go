package member

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

type UpdateParams struct {
	Position *string
	Label    *string
}

func (s Service) UpdateMember(ctx context.Context, ID uuid.UUID, params UpdateParams) (member entity.Member, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		member, err = s.repo.UpdateMember(ctx, ID, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update member %s: %w", ID, err),
			)
		}

		if err = s.messenger.WriteMemberUpdated(ctx, member); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send member updated message for member %s: %w", ID, err),
			)
		}

		return nil
	}); err != nil {
		return entity.Member{}, err
	}

	return member, nil
}

func (s Service) UpdateMemberByUser(
	ctx context.Context,
	accountID, memberID uuid.UUID,
	params UpdateParams,
) (entity.Member, error) {
	member, err := s.GetMember(ctx, memberID)
	if err != nil {
		return entity.Member{}, err
	}

	initiator, err := s.GetInitiatorMember(ctx, accountID, memberID)
	if err != nil {
		return entity.Member{}, err
	}

	err = s.CheckAccessToManageOtherMember(ctx, initiator.ID, member.ID)
	if err != nil {
		return entity.Member{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member %s has no permission to manage members: %w", initiator.ID, err),
		)
	}

	return s.UpdateMember(ctx, member.ID, params)
}

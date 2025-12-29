package member

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

func (s Service) DeleteMember(ctx context.Context, ID uuid.UUID) error {
	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		err := s.repo.DeleteMember(ctx, ID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete member %s: %w", ID, err),
			)
		}

		if err = s.messenger.WriteMemberDeleted(ctx, ID); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send member deleted message for member %s: %w", ID, err),
			)
		}

		return nil
	})
}

func (s Service) DeleteMemberByUser(ctx context.Context, accountID, memberID uuid.UUID) error {
	member, err := s.GetMember(ctx, memberID)
	if err != nil {
		return err
	}

	initiator, err := s.GetInitiatorMember(ctx, accountID, member.AgglomerationID)
	if err != nil {
		return err
	}

	err = s.CheckAccessToManageOtherMember(ctx, initiator.ID, member.ID)
	if err != nil {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member %s has no permission to manage members: %w", initiator.ID, err),
		)
	}

	return s.DeleteMember(ctx, member.ID)
}

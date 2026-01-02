package member

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/domain/errx"
)

func (s Service) DeleteMember(ctx context.Context, accountID, memberID uuid.UUID) error {
	member, err := s.GetMemberByID(ctx, memberID)
	if err != nil {
		return err
	}

	initiator, err := s.GetInitiatorMember(ctx, accountID, member.OrganizationID)
	if err != nil {
		return err
	}

	err = s.CheckAccessToManageOtherMember(ctx, initiator.ID, member.ID)
	if err != nil {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member %s has no permission to manage members: %w", initiator.ID, err),
		)
	}

	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		err = s.repo.DeleteMember(ctx, memberID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete member %s: %w", memberID, err),
			)
		}

		if err = s.messenger.WriteMemberDeleted(ctx, member); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send member deleted message for member %s: %w", memberID, err),
			)
		}

		return nil
	})
}

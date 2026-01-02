package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (s Service) DeclineInvite(
	ctx context.Context,
	accountID, inviteID uuid.UUID,
) (invite models.Invite, err error) {
	invite, err = s.getInvite(ctx, accountID)
	if err != nil {
		return models.Invite{}, err
	}

	if invite.AccountID != accountID {
		return models.Invite{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("account has no rights to decline this invite"),
		)
	}
	if invite.Status != models.InviteStatusSent {
		return models.Invite{}, errx.ErrorInviteAlreadyAnswered.Raise(
			fmt.Errorf("invite status is %s", invite.Status),
		)
	}
	if invite.ExpiresAt.Before(time.Now().UTC()) {
		return models.Invite{}, errx.ErrorInviteExpired.Raise(
			fmt.Errorf("invite expired at %s", invite.ExpiresAt),
		)
	}

	err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		invite, err = s.repo.UpdateInviteStatus(ctx, inviteID, models.InviteStatusDeclined)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update invite status: %w", err),
			)
		}

		err = s.messenger.WriteInviteDeclined(ctx, invite)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to write declined invite event: %w", err),
			)
		}

		return nil
	})

	return invite, err
}

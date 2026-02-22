package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
	"github.com/netbill/organizations-svc/internal/core/errx"
)

func (m *Module) Decline(
	ctx context.Context,
	initiator domain.AccountActor,
	inviteID uuid.UUID,
) (invite domain.Invite, err error) {
	invite, err = m.GetForAccount(ctx, initiator, inviteID)
	if err != nil {
		return domain.Invite{}, err
	}

	if invite.AccountID != initiator {
		return domain.Invite{}, errx.ErrorInviteNotForInitiator.Raise(
			fmt.Errorf("account has no rights to decline this invite"),
		)
	}
	if invite.Status != domain.InviteStatusSent {
		return domain.Invite{}, errx.ErrorInviteAlreadyAnswered.Raise(
			fmt.Errorf("invite status is %s", invite.Status),
		)
	}
	if invite.ExpiresAt.Before(time.Now().UTC()) {
		return domain.Invite{}, errx.ErrorInviteExpired.Raise(
			fmt.Errorf("invite expired at %s", invite.ExpiresAt),
		)
	}

	err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		invite, err = m.repo.UpdateInviteStatus(ctx, inviteID, domain.InviteStatusDeclined)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrgInviteDeclined(ctx, invite)
		if err != nil {
			return err
		}

		return nil
	})

	return invite, err
}

package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) Delete(
	ctx context.Context,
	initiator models.Initiator,
	inviteID uuid.UUID,
) error {
	invite, err := m.GetForAccount(ctx, initiator.GetAccountID(), inviteID)
	if err != nil {
		return err
	}

	if invite.AccountID != initiator.GetAccountID() {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("account has no rights to accept this invite"),
		)
	}
	if invite.Status != models.InviteStatusSent {
		return errx.ErrorInviteAlreadyAnswered.Raise(
			fmt.Errorf("invite status is %s", invite.Status),
		)
	}
	if invite.ExpiresAt.Before(time.Now().UTC()) {
		return errx.ErrorInviteExpired.Raise(
			fmt.Errorf("invite expired at %s", invite.ExpiresAt),
		)
	}

	if err = m.checkPermissionForManageInvites(
		ctx,
		initiator,
		invite.OrganizationID,
	); err != nil {
		return err
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		err = m.repo.DeleteInvite(ctx, inviteID)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrgInviteDeleted(ctx, invite)
		if err != nil {
			return err
		}

		return nil
	})
}

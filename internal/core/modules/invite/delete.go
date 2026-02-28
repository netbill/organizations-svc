package invite

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) Delete(
	ctx context.Context,
	actor models.AccountActor,
	inviteID uuid.UUID,
) error {
	invite, err := m.repo.GetInvite(ctx, inviteID)
	if errors.Is(err, errx.ErrorInviteNotFound) {
		buried, err := m.repo.InviteIsBuried(ctx, inviteID)
		if err != nil {
			return err
		}
		if buried {
			return errx.ErrorInviteDeleted.Raise(
				fmt.Errorf("invite with id %s is already deleted", inviteID),
			)
		}
	}
	if err != nil {
		return err
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

	initiator, err := m.getInitiator(ctx, actor, invite.OrganizationID)
	if err != nil {
		return err
	}
	if !initiator.Head {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("account has no rights to create invite for this organization"),
		)
	}

	org, err := m.repo.GetOrganizationByID(ctx, invite.OrganizationID)
	if err != nil {
		return err
	}
	if org.Status == models.OrganizationStatusSuspended {
		return errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization with id %s is suspended", invite.OrganizationID),
		)
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		err = m.repo.BuryInvite(ctx, inviteID)
		if err != nil {
			return err
		}

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

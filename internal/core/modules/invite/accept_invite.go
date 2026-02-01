package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) AcceptInvite(
	ctx context.Context,
	accountID, inviteID uuid.UUID,
) (invite models.Invite, err error) {
	invite, err = m.GetInviteForAccount(ctx, accountID, inviteID)
	if err != nil {
		return models.Invite{}, err
	}

	if invite.AccountID != accountID {
		return models.Invite{}, errx.ErrorInviteNotForInitiator.Raise(
			fmt.Errorf("account has no rights to accept this invite"),
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

	if _, err = m.checkOrganizationIsActiveAndExists(ctx, invite.OrganizationID); err != nil {
		return models.Invite{}, err
	}

	err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		invite, err = m.repo.UpdateInviteStatus(ctx, inviteID, models.InviteStatusAccepted)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrgInviteAccepted(ctx, invite)
		if err != nil {
			return err
		}

		mem, err := m.repo.CreateMember(ctx, accountID, invite.OrganizationID)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrgMemberCreated(ctx, mem)
		if err != nil {
			return err
		}

		return nil
	})

	return invite, err
}

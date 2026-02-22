package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
	"github.com/netbill/organizations-svc/internal/core/errx"
)

func (m *Module) Accept(
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
			fmt.Errorf("account has no rights to accept this invite"),
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

	if _, err = m.checkOrganizationIsActiveAndExists(ctx, invite.OrganizationID); err != nil {
		return domain.Invite{}, err
	}

	err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		invite, err = m.repo.UpdateInviteStatus(ctx, inviteID, domain.InviteStatusAccepted)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrgInviteAccepted(ctx, invite)
		if err != nil {
			return err
		}

		mem, err := m.repo.CreateMember(ctx, initiator, invite.OrganizationID)
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

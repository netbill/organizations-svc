package invite

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Invite) Cancelled(
	ctx context.Context,
	actor models.AccountActor,
	inviteID uuid.UUID,
) (models.Invite, error) {
	invite, err := m.repo.GetInvite(ctx, inviteID)
	if err != nil {
		return models.Invite{}, err
	}

	initiator, err := m.getInitiator(ctx, actor, invite.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}
	if !initiator.Head {
		return models.Invite{}, errx.ErrorNotOrganizationHead.Raise(
			fmt.Errorf("account has no rights to create invite for this organization"),
		)
	}

	if invite.Status == models.InviteStatusCancelled {
		return models.Invite{}, nil
	} else if invite.Status != models.InviteStatusSent {
		return models.Invite{}, errx.ErrorInviteAlreadyAnswered.Raise(
			fmt.Errorf("invite status is %s", invite.Status),
		)
	}

	org, err := m.repo.GetOrganizationByID(ctx, invite.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}
	if org.Status == models.OrganizationStatusSuspended {
		return models.Invite{}, errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization with id %s is suspended", invite.OrganizationID),
		)
	}

	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		invite, err = m.repo.UpdateInviteStatus(ctx, inviteID, models.InviteStatusCancelled)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrgInviteCanceled(ctx, invite)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Invite{}, err
	}

	return invite, nil

}

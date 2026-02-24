package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) Accept(
	ctx context.Context,
	actor models.AccountActor,
	inviteID uuid.UUID,
) (invite models.Invite, err error) {
	invite, err = m.repo.GetInvite(ctx, inviteID)
	if err != nil {
		return models.Invite{}, err
	}

	if invite.AccountID != actor {
		return models.Invite{}, errx.ErrorInviteNotForInitiator.Raise(
			fmt.Errorf("account has no rights to accept this invite"),
		)
	}
	if invite.Status == models.InviteStatusAccepted {
		return invite, nil
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

	org, err := m.repo.GetOrganizationByID(ctx, invite.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}
	if org.Status != models.OrganizationStatusActive {
		return models.Invite{}, errx.ErrorOrganizationIsNotActive.Raise(
			fmt.Errorf("organization with id %s is not active", invite.OrganizationID),
		)
	}

	err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		invite, err = m.repo.UpdateInviteStatus(ctx, inviteID, models.InviteStatusAccepted)
		if err != nil {
			return err
		}

		mem, err := m.repo.CreateMember(ctx, actor, invite.OrganizationID)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrgInviteAccepted(ctx, invite)
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

func (m *Module) Decline(
	ctx context.Context,
	actor models.AccountActor,
	inviteID uuid.UUID,
) (invite models.Invite, err error) {
	invite, err = m.repo.GetInvite(ctx, inviteID)
	if err != nil {
		return models.Invite{}, err
	}

	if invite.AccountID != actor {
		return models.Invite{}, errx.ErrorInviteNotForInitiator.Raise(
			fmt.Errorf("account has no rights to decline this invite"),
		)
	}
	if invite.Status == models.InviteStatusDeclined {
		return invite, nil
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

	org, err := m.repo.GetOrganizationByID(ctx, invite.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}
	if org.Status != models.OrganizationStatusActive {
		return models.Invite{}, errx.ErrorOrganizationIsNotActive.Raise(
			fmt.Errorf("organization with id %s is not active", invite.OrganizationID),
		)
	}

	err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		invite, err = m.repo.UpdateInviteStatus(ctx, inviteID, models.InviteStatusDeclined)
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

package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/errx"
	"github.com/netbill/organizations-svc/internal/models"
)

func (m *Service) Accept(
	ctx context.Context,
	actor models.AccountActor,
	inviteID uuid.UUID,
) (invite models.Invite, err error) {
	invite, err = m.validateInviteForResponse(ctx, actor, inviteID, models.InviteStatusAccepted)
	if err != nil {
		return models.Invite{}, err
	}
	if invite.Status == models.InviteStatusAccepted {
		return invite, nil
	}

	_, err = m.auth.ValidateOrg(ctx, invite.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}

	if err = m.tx.Transaction(ctx, func(ctx context.Context) error {
		invite, err = m.repo.UpdateStatus(ctx, inviteID, models.InviteStatusAccepted)
		if err != nil {
			return err
		}

		mem, err := m.member.Create(ctx, actor, invite.OrganizationID, false)
		if err != nil {
			return err
		}

		if err = m.messenger.WriteOrgInviteAccepted(ctx, invite); err != nil {
			return err
		}

		return m.messenger.WriteOrgMemberCreated(ctx, mem)
	}); err != nil {
		return models.Invite{}, err
	}

	return invite, nil
}

func (m *Service) Decline(
	ctx context.Context,
	actor models.AccountActor,
	inviteID uuid.UUID,
) (invite models.Invite, err error) {
	invite, err = m.validateInviteForResponse(ctx, actor, inviteID, models.InviteStatusDeclined)
	if err != nil {
		return models.Invite{}, err
	}
	if invite.Status == models.InviteStatusDeclined {
		return invite, nil
	}

	if err = m.tx.Transaction(ctx, func(ctx context.Context) error {
		invite, err = m.repo.UpdateStatus(ctx, inviteID, models.InviteStatusDeclined)
		if err != nil {
			return err
		}

		return m.messenger.WriteOrgInviteDeclined(ctx, invite)
	}); err != nil {
		return models.Invite{}, err
	}

	return invite, nil
}

func (m *Service) Cancelled(
	ctx context.Context,
	actor models.AccountActor,
	inviteID uuid.UUID,
) (models.Invite, error) {
	invite, err := m.repo.Get(ctx, inviteID)
	if err != nil {
		return models.Invite{}, err
	}

	_, err = m.auth.ValidateOrg(ctx, invite.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}

	_, err = m.auth.AuthorizeOrgHead(ctx, actor, invite.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}

	if invite.Status == models.InviteStatusCancelled {
		return invite, nil
	}
	if invite.Status != models.InviteStatusSent {
		return models.Invite{}, errx.ErrorInviteAlreadyAnswered.Raise(
			fmt.Errorf("invite status is %s", invite.Status),
		)
	}

	if err = m.tx.Transaction(ctx, func(ctx context.Context) error {
		invite, err = m.repo.UpdateStatus(ctx, inviteID, models.InviteStatusCancelled)
		if err != nil {
			return err
		}

		return m.messenger.WriteOrgInviteCanceled(ctx, invite)
	}); err != nil {
		return models.Invite{}, err
	}

	return invite, nil
}

func (m *Service) validateInviteForResponse(
	ctx context.Context,
	actor models.AccountActor,
	inviteID uuid.UUID,
	terminalStatus string,
) (models.Invite, error) {
	invite, err := m.repo.Get(ctx, inviteID)
	if err != nil {
		return models.Invite{}, err
	}

	if invite.AccountID != actor {
		return models.Invite{}, errx.ErrorInviteNotForInitiator.Raise(
			fmt.Errorf("account has no rights to respond to this invite"),
		)
	}
	if invite.Status == terminalStatus {
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

	return invite, nil
}

package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/domain/errx"
	"github.com/netbill/organizations-svc/internal/domain/models"
)

func (s Service) AcceptInvite(
	ctx context.Context,
	accountID, inviteID uuid.UUID,
) (invite models.Invite, err error) {
	invite, err = s.getInvite(ctx, accountID)
	if err != nil {
		return models.Invite{}, err
	}

	if invite.AccountID != accountID {
		return models.Invite{}, errx.ErrorNotEnoughRights.Raise(
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

	if _, err = s.checkOrganizationIsActiveAndExists(ctx, invite.OrganizationID); err != nil {
		return models.Invite{}, err
	}

	err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		invite, err = s.repo.UpdateInviteStatus(ctx, inviteID, models.InviteStatusAccepted)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update invite status: %w", err),
			)
		}

		err = s.messenger.WriteInviteAccepted(ctx, invite)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to write accepted invite event: %w", err),
			)
		}

		mem, err := s.repo.CreateMember(ctx, accountID, invite.ID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create member from invite: %w", err),
			)
		}

		err = s.messenger.WriteMemberCreated(ctx, mem)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to write created member event: %w", err),
			)
		}

		return nil
	})

	return invite, err
}

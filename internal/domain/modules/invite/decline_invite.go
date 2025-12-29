package invite

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
)

func (s Service) DeclineInvite(
	ctx context.Context,
	accountID, inviteID uuid.UUID,
) (invite models.Invite, err error) {
	invite, err = s.repo.GetInviteByID(ctx, inviteID)
	if err != nil {
		return models.Invite{}, err
	}

	if invite.AccountID != accountID {
		return models.Invite{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("account has no rights to decline this invite"),
		)
	}
	if invite.Status != models.InviteStatusSent {
		return models.Invite{}, err
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		invite, err = s.repo.UpdateInviteStatus(ctx, inviteID, models.InviteStatusDeclined)
		if err != nil {
			return err
		}

		err = s.messenger.WriteDeclinedInvite(ctx, invite)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Invite{}, err
	}

	return invite, nil
}

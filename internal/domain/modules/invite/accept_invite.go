package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
)

func (s Service) AcceptInvite(
	ctx context.Context,
	accountID, inviteID uuid.UUID,
) (invite models.Invite, err error) {
	invite, err = s.repo.GetInviteByID(ctx, inviteID)
	if err != nil {
		return models.Invite{}, err
	}

	if invite.AccountID != accountID {
		return models.Invite{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("account has no rights to accept this invite"),
		)
	}

	if invite.Status != models.InviteStatusSent {
		return models.Invite{}, err
	}
	if invite.ExpiresAt.Before(time.Now().UTC()) {
		return models.Invite{}, err
	}

	err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		invite, err = s.repo.UpdateInviteStatus(ctx, inviteID, models.InviteStatusAccepted)
		if err != nil {
			return err
		}

		err = s.messenger.WriteAcceptedInvite(ctx, invite)
		if err != nil {
			return err
		}

		mem, err := s.repo.CreateMember(ctx, accountID, invite.ID)
		if err != nil {
			return err
		}

		err = s.messenger.WriteCreatedNewMember(ctx, mem)
		if err != nil {
			return err
		}

		return nil
	})

	return invite, nil
}

package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

func (s Service) AcceptInvite(
	ctx context.Context,
	accountID, inviteID uuid.UUID,
) (invite entity.Invite, err error) {
	invite, err = s.repo.GetInviteByID(ctx, inviteID)
	if err != nil {
		return entity.Invite{}, err
	}

	if invite.AccountID != accountID {
		return entity.Invite{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("account has no rights to accept this invite"),
		)
	}

	if invite.Status != entity.InviteStatusSent {
		return entity.Invite{}, err
	}
	if invite.ExpiresAt.Before(time.Now().UTC()) {
		return entity.Invite{}, err
	}

	err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		invite, err = s.repo.UpdateInviteStatus(ctx, inviteID, entity.InviteStatusAccepted)
		if err != nil {
			return err
		}

		_, err = s.repo.CreateMember(ctx, accountID, invite.ID)
		if err != nil {
			return err
		}

		return nil
	})

	return invite, nil
}

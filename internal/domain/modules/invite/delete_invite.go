package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
)

func (s Service) DeleteInvite(
	ctx context.Context,
	accountID, inviteID uuid.UUID,
) error {
	invite, err := s.repo.GetInviteByID(ctx, inviteID)
	if err != nil {
		return err
	}

	if invite.AccountID != accountID {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("account has no rights to accept this invite"),
		)
	}

	if invite.Status != models.InviteStatusSent {
		return err
	}
	if invite.ExpiresAt.Before(time.Now().UTC()) {
		return err
	}

	if err = s.checkPermissionForManageInvite(
		ctx,
		accountID,
		invite.AgglomerationID,
	); err != nil {
		return err
	}

	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		err = s.repo.DeleteInvite(ctx, inviteID)
		if err != nil {
			return err
		}

		err = s.messenger.WriteDeletedInvite(ctx, invite.ID)
		if err != nil {
			return err
		}

		return nil
	})
}

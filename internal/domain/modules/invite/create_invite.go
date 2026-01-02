package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
)

type CreateParams struct {
	AccountID       uuid.UUID
	AgglomerationID uuid.UUID
	ExpiresAt       time.Time
}

func (s Service) CreateInvite(
	ctx context.Context,
	accountID uuid.UUID,
	params CreateParams,
) (invite models.Invite, err error) {
	initiator, err := s.getInitiator(ctx, accountID, params.AgglomerationID)
	if err != nil {
		return invite, err
	}

	if err = s.checkPermissionForManageInvite(
		ctx,
		initiator.ID,
	); err != nil {
		return models.Invite{}, err
	}

	if _, err = s.checkAgglomerationIsActiveAndExists(ctx, params.AgglomerationID); err != nil {
		return models.Invite{}, err
	}

	err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		invite, err = s.repo.CreateInvite(ctx, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create invite: %w", err),
			)
		}

		err = s.messenger.WriteInviteCreated(ctx, invite)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to write created invite event: %w", err),
			)
		}

		return nil
	})

	return invite, err
}

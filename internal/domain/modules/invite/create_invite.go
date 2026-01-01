package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
)

type CreateParams struct {
	AccountID       uuid.UUID
	AgglomerationID uuid.UUID
	ExpiresAt       time.Time
}

func (s Service) CreateInvite(ctx context.Context, params CreateParams) (invite models.Invite, err error) {
	_, err = s.checkAgglomerationIsActiveAndExists(ctx, params.AgglomerationID)
	if err != nil {
		return models.Invite{}, err
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
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
	}); err != nil {
		return models.Invite{}, err
	}

	return invite, nil
}

func (s Service) SentInviteByUser(
	ctx context.Context,
	accountID uuid.UUID,
	params CreateParams,
) (models.Invite, error) {
	if err := s.checkPermissionForManageInvite(
		ctx,
		accountID,
		params.AgglomerationID,
	); err != nil {
		return models.Invite{}, err
	}

	res, err := s.CreateInvite(ctx, params)
	if err != nil {
		return models.Invite{}, err
	}

	return res, nil
}

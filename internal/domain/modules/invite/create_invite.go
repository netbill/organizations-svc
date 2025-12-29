package invite

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
)

type CreateParams struct {
	AccountID       uuid.UUID
	AgglomerationID uuid.UUID
	ExpiresAt       time.Time
}

func (s Service) CreateInvite(ctx context.Context, params CreateParams) (entity.Invite, error) {
	res, err := s.repo.CreateInvite(ctx, params)
	if err != nil {
		return entity.Invite{}, err
	}

	return res, nil
}

func (s Service) CreateInviteByUser(
	ctx context.Context,
	accountID uuid.UUID,
	params CreateParams,
) (entity.Invite, error) {
	if err := s.checkPermissionForManageInvite(
		ctx,
		accountID,
		params.AgglomerationID,
	); err != nil {
		return entity.Invite{}, err
	}

	res, err := s.CreateInvite(ctx, params)
	if err != nil {
		return entity.Invite{}, err
	}

	return res, nil
}

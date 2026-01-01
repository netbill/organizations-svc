package invite

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
)

func (s Service) GetInvite(ctx context.Context, id uuid.UUID) (models.Invite, error) {
	res, err := s.repo.GetInviteByID(ctx, id)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get invite by id %s: %w", id.String(), err),
		)
	}
	if res.IsNil() {
		return models.Invite{}, errx.ErrorInviteNotFound.Raise(
			fmt.Errorf("invite with id %s not found", id.String()),
		)
	}

	return res, nil
}

type FilterInviteParams struct {
	AgglomerationID *uuid.UUID
	AccountID       *uuid.UUID
	Status          *string
}

func (s Service) FilterInvites(
	ctx context.Context,
	filter FilterInviteParams,
) ([]models.Invite, error) {
	res, err := s.repo.FilterInvites(ctx, filter)
	if err != nil {
		return []models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to filter invites: %w", err),
		)
	}

	return res, nil
}

package invite

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/models"
)

func (s Service) GetInvite(ctx context.Context, id uuid.UUID) (models.Invite, error) {
	res, err := s.repo.GetInviteByID(ctx, id)
	if err != nil {
		return models.Invite{}, err
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
		return []models.Invite{}, err
	}

	return res, nil
}

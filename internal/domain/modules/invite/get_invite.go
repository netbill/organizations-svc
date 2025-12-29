package invite

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
)

func (s Service) GetInvite(ctx context.Context, id uuid.UUID) (entity.Invite, error) {
	res, err := s.repo.GetInviteByID(ctx, id)
	if err != nil {
		return entity.Invite{}, err
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
) ([]entity.Invite, error) {
	res, err := s.repo.FilterInvites(ctx, filter)
	if err != nil {
		return []entity.Invite{}, err
	}

	return res, nil
}

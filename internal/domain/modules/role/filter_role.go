package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/models"
	"github.com/umisto/pagi"
)

type FilterParams struct {
	AgglomerationID *uuid.UUID
	RolesID         *[]uuid.UUID
	Head            *bool
	Rank            *int
	Name            *string
}

func (s Service) FilterRoles(
	ctx context.Context,
	params FilterParams,
	offset uint,
	limit uint,
) (pagi.Page[[]models.Role], error) {
	return s.repo.FilterRoles(ctx, params, offset, limit)
}

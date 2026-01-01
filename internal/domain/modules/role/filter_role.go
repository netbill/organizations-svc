package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
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
	res, err := s.repo.FilterRoles(ctx, params, offset, limit)
	if err != nil {
		return pagi.Page[[]models.Role]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to filter roles: %w", err),
		)
	}

	return res, nil
}

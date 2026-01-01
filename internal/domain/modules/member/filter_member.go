package member

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
	AccountID       *uuid.UUID
	RoleID          *uuid.UUID
	Username        *string
	BestMatch       *string
	PermissionCode  *string
	Label           *string
	Position        *string
	RoleRankUp      *uint
	RoleRankDown    *uint
}

func (s Service) FilterMembers(
	ctx context.Context,
	filter FilterParams,
	offset uint,
	limit uint,
) (pagi.Page[[]models.Member], error) {
	res, err := s.repo.FilterMembers(ctx, filter, offset, limit)
	if err != nil {
		return pagi.Page[[]models.Member]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to filter members: %w", err),
		)
	}

	return res, nil
}

package member

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/models"
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
	return s.repo.FilterMembers(ctx, filter, offset, limit)
}

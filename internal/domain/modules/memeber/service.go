package memeber

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/pagi"
)

type repo interface {
	CreateMember(ctx context.Context, accountID, agglomerationID uuid.UUID) error

	GetMember(ctx context.Context, memberID uuid.UUID) (entity.Member, error)
	GetMemberByAccountAndAgglomeration(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
	) (entity.Member, error)

	UpdateMember(ctx context.Context, memberID uuid.UUID, params UpdateParams) error

	DeleteMember(ctx context.Context, memberID uuid.UUID) error

	FilterMembers(
		ctx context.Context,
		filter FilterParams,
		pagination pagi.Params,
	) (pagi.Page[entity.Member], error)
}

type FilterParams struct {
	AgglomerationID *uuid.UUID
	AccountID       *uuid.UUID
	Username        *string
	RoleID          *uuid.UUID
	PermissionCode  *string
}

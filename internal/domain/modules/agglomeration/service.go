package agglomeration

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/pagi"
)

type Service struct {
	repo     repo
	messager messager
}

type repo interface {
	CreateAgglomeration(ctx context.Context, name string) (entity.Agglomeration, error)
	FilterAgglomerations(
		ctx context.Context,
		filter FilterParams,
		pagination pagi.Params,
	) (pagi.Page[entity.Agglomeration], error)

	UpdateAgglomeration(
		ctx context.Context,
		ID uuid.UUID,
		params UpdateParams,
	) (entity.Agglomeration, error)

	UpdateAgglomerationStatus(ctx context.Context, ID uuid.UUID, status string) (entity.Agglomeration, error)

	GetAgglomerationByID(ctx context.Context, ID uuid.UUID) (entity.Agglomeration, error)

	DeleteAgglomeration(ctx context.Context, ID uuid.UUID) error

	CheckMemberHavePermissionByCode(ctx context.Context, memberID uuid.UUID, permissionKey string) (bool, error)
	CheckMemberHavePermissionByID(ctx context.Context, memberID, permissionID uuid.UUID) (bool, error)
	CheckAccountHavePermissionByCode(ctx context.Context, accountID uuid.UUID, permissionKey string) (bool, error)
	CheckAccountHavePermissionByID(ctx context.Context, accountID, permissionID uuid.UUID) (bool, error)
}

type messager interface {
	WriteAgglomerationCreated(ctx context.Context, agglomeration entity.Agglomeration) error

	WriteAgglomerationActivated(ctx context.Context, agglomeration entity.Agglomeration) error
	WriteAgglomerationDeactivated(ctx context.Context, agglomeration entity.Agglomeration) error
	WriteAgglomerationSuspended(ctx context.Context, agglomeration entity.Agglomeration) error

	WriteAgglomerationUpdated(ctx context.Context, agglomeration entity.Agglomeration) error

	WriteAgglomerationDeleted(ctx context.Context, agglomerationID uuid.UUID) error
}

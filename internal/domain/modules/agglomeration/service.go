package agglomeration

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/pagi"
)

type Service struct {
	repo Repo
}

type Repo interface {
	CreateAgglomeration(ctx context.Context, name string) (entity.Agglomeration, error)
	UpdateAgglomeration(ctx context.Context, ID uuid.UUID, params UpdateParams) (entity.Agglomeration, error)
	ActivateAgglomeration(ctx context.Context, ID uuid.UUID) error
	DeactivateAgglomeration(ctx context.Context, ID uuid.UUID) error

	GetAgglomerationByID(ctx context.Context, ID uuid.UUID) (entity.Agglomeration, error)
	FilterAgglomerations(
		ctx context.Context,
		params FilterParams,
		pagination pagi.Params,
	) (pagi.Page[entity.Agglomeration], error)

	DeleteAgglomeration(ctx context.Context, ID uuid.UUID) error
}

package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/pagi"
)

type Service struct {
	repo      repo
	messenger messanger
}

type repo interface {
	CreateAgglomeration(ctx context.Context, name string) (entity.Agglomeration, error)

	UpdateAgglomeration(
		ctx context.Context,
		ID uuid.UUID,
		params UpdateParams,
	) (entity.Agglomeration, error)
	UpdateAgglomerationStatus(ctx context.Context, ID uuid.UUID, status string) (entity.Agglomeration, error)

	GetAgglomerationByID(ctx context.Context, ID uuid.UUID) (entity.Agglomeration, error)

	DeleteAgglomeration(ctx context.Context, ID uuid.UUID) error

	FilterAgglomerations(
		ctx context.Context,
		filter FilterParams,
		offset, limit uint,
	) (pagi.Page[[]entity.Agglomeration], error)

	CheckAccountHavePermissionByCode(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
		permissionKey string,
	) (bool, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messanger interface {
	WriteAgglomerationCreated(ctx context.Context, agglomeration entity.Agglomeration) error

	WriteAgglomerationActivated(ctx context.Context, agglomeration entity.Agglomeration) error
	WriteAgglomerationDeactivated(ctx context.Context, agglomeration entity.Agglomeration) error
	WriteAgglomerationSuspended(ctx context.Context, agglomeration entity.Agglomeration) error

	WriteAgglomerationUpdated(ctx context.Context, agglomeration entity.Agglomeration) error

	WriteAgglomerationDeleted(ctx context.Context, agglomerationID uuid.UUID) error
}

func (s Service) checkPermissionForManageAgglomeration(
	ctx context.Context,
	accountID uuid.UUID,
	agglomerationID uuid.UUID,
) error {
	access, err := s.repo.CheckAccountHavePermissionByCode(
		ctx,
		accountID,
		agglomerationID,
		entity.RolePermissionManageAgglomeration.String(),
	)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check initiator permissions: %w", err))
	}
	if !access {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator has no access to activate agglomeration"),
		)
	}

	return nil
}

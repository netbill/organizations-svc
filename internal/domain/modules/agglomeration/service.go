package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
	"github.com/umisto/pagi"
)

type Service struct {
	repo      repo
	messenger messanger
}

func New(repo repo, messenger messanger) Service {
	return Service{
		repo:      repo,
		messenger: messenger,
	}
}

type repo interface {
	CreateAgglomeration(ctx context.Context, params CreateParams) (models.Agglomeration, error)

	UpdateAgglomeration(
		ctx context.Context,
		ID uuid.UUID,
		params UpdateParams,
	) (models.Agglomeration, error)
	UpdateAgglomerationStatus(ctx context.Context, ID uuid.UUID, status string) (models.Agglomeration, error)

	GetAgglomerationByID(ctx context.Context, ID uuid.UUID) (models.Agglomeration, error)

	DeleteAgglomeration(ctx context.Context, ID uuid.UUID) error

	FilterAgglomerations(
		ctx context.Context,
		filter FilterParams,
		offset, limit uint,
	) (pagi.Page[[]models.Agglomeration], error)

	CheckAccountHavePermissionByCode(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
		permissionKey string,
	) (bool, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messanger interface {
	WriteAgglomerationCreated(ctx context.Context, agglomeration models.Agglomeration) error

	WriteAgglomerationActivated(ctx context.Context, agglomeration models.Agglomeration) error
	WriteAgglomerationDeactivated(ctx context.Context, agglomeration models.Agglomeration) error
	WriteAgglomerationSuspended(ctx context.Context, agglomeration models.Agglomeration) error

	WriteAgglomerationUpdated(ctx context.Context, agglomeration models.Agglomeration) error

	WriteAgglomerationDeleted(ctx context.Context, agglomeration models.Agglomeration) error
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
		models.RolePermissionManageAgglomeration.String(),
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

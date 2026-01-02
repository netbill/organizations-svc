package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
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

	GetAgglomerationByID(ctx context.Context, ID uuid.UUID) (models.Agglomeration, error)
	GetAgglomerations(
		ctx context.Context,
		filter FilterParams,
		offset, limit uint,
	) (pagi.Page[[]models.Agglomeration], error)
	GetAgglomerationsForUser(
		ctx context.Context,
		accountID uuid.UUID,
		limit, offset uint,
	) (pagi.Page[[]models.Agglomeration], error)

	UpdateAgglomeration(
		ctx context.Context,
		ID uuid.UUID,
		params UpdateParams,
	) (models.Agglomeration, error)
	UpdateAgglomerationStatus(ctx context.Context, ID uuid.UUID, status string) (models.Agglomeration, error)
	UpdateAgglomerationMaxRoles(ctx context.Context, ID uuid.UUID, maxRoles uint) (models.Agglomeration, error)

	DeleteAgglomeration(ctx context.Context, ID uuid.UUID) error

	CheckMemberHavePermission(
		ctx context.Context,
		memberID uuid.UUID,
		permissionCode string,
	) (bool, error)
	GetMemberByAccountAndAgglomeration(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
	) (models.Member, error)

	GetMemberMaxRole(ctx context.Context, memberID uuid.UUID) (models.Role, error)

	CreateMember(ctx context.Context, accountID, agglomerationID uuid.UUID) (models.Member, error)
	CreateHeadRole(ctx context.Context, agglomerationID uuid.UUID) (models.Role, error)
	AddMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messanger interface {
	WriteAgglomerationCreated(ctx context.Context, agglomeration models.Agglomeration) error

	WriteAgglomerationActivated(ctx context.Context, agglomeration models.Agglomeration) error
	WriteAgglomerationDeactivated(ctx context.Context, agglomeration models.Agglomeration) error
	WriteAgglomerationSuspended(ctx context.Context, agglomeration models.Agglomeration) error

	WriteAgglomerationUpdated(ctx context.Context, agglomeration models.Agglomeration) error

	WriteAgglomerationDeleted(ctx context.Context, agglomeration models.Agglomeration) error
	WriteRoleCreated(ctx context.Context, role models.Role) error
}

func (s Service) chekPermissionForManageAgglomeration(
	ctx context.Context,
	memberID uuid.UUID,
) error {
	access, err := s.repo.CheckMemberHavePermission(
		ctx,
		memberID,
		models.RolePermissionManageAgglomeration,
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

func (s Service) getInitiator(ctx context.Context, accountID, agglomerationID uuid.UUID) (models.Member, error) {
	row, err := s.repo.GetMemberByAccountAndAgglomeration(ctx, accountID, agglomerationID)
	if err != nil {
		return models.Member{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get member with account id %s and agglomeration id %s: %w",
				accountID, agglomerationID, err),
		)
	}
	if row.IsNil() {
		return models.Member{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member with account id %s and agglomeration id %s not found", accountID, agglomerationID),
		)
	}

	return row, nil
}

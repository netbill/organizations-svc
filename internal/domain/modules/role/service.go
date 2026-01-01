package role

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
	messenger messenger
}

func New(repo repo, messenger messenger) Service {
	return Service{
		repo:      repo,
		messenger: messenger,
	}
}

type repo interface {
	CreateRole(ctx context.Context, params CreateParams) (models.Role, error)

	UpdateRole(ctx context.Context, roleID uuid.UUID, params UpdateParams) (models.Role, error)
	UpdateRoleRank(ctx context.Context, roleID uuid.UUID, newRank uint) (models.Role, error)
	UpdateRolesRanks(
		ctx context.Context,
		agglomerationID uuid.UUID,
		order map[uuid.UUID]uint,
	) error

	DeleteRole(ctx context.Context, roleID uuid.UUID) error

	FilterRoles(
		ctx context.Context,
		filter FilterParams,
		offset uint,
		limit uint,
	) (pagi.Page[[]models.Role], error)
	GetRole(ctx context.Context, roleID uuid.UUID) (models.Role, error)

	GetMember(ctx context.Context, memberID uuid.UUID) (models.Member, error)
	GetMemberByAccountAndAgglomeration(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
	) (models.Member, error)

	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]models.Permission, error)
	SetRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
		permissions map[models.CodeRolePermission]bool,
	) error

	GetAllPermissions(ctx context.Context) ([]models.Permission, error)

	GetAccountMaxRoleInAgglomeration(ctx context.Context, accountID, agglomerationID uuid.UUID) (models.Role, error)
	GetMemberMaxRole(ctx context.Context, memberID uuid.UUID) (models.Role, error)

	CanInteract(ctx context.Context, firstMemberID, secondMemberID uuid.UUID) (bool, error)

	GetMemberRoles(ctx context.Context, memberID uuid.UUID) ([]models.Role, error)
	DeleteMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error
	AddMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error

	CheckAccountHavePermissionByCode(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
		permissionKey string,
	) (bool, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteRoleCreated(ctx context.Context, role models.Role) error
	WriteRoleUpdated(ctx context.Context, role models.Role) error
	WriteRoleDeleted(ctx context.Context, role models.Role) error

	WriteRolesRanksUpdated(
		ctx context.Context,
		agglomerationID uuid.UUID,
		order map[uuid.UUID]uint,
	) error
	WriteRolePermissionsUpdated(
		ctx context.Context,
		roleID uuid.UUID,
		permissions map[models.CodeRolePermission]bool,
	) error
}

func (s Service) GetInitiator(ctx context.Context, memberID uuid.UUID) (models.Member, error) {
	row, err := s.repo.GetMember(ctx, memberID)
	if err != nil {
		return models.Member{}, err
	}
	if row.IsNil() {
		return models.Member{}, errx.ErrorMemberNotFound.Raise(
			fmt.Errorf("member with id %s not found", memberID),
		)
	}

	return row, nil
}

func (s Service) GetInitiatorByAccountAndAgglomeration(
	ctx context.Context,
	initiatorAccountID, agglomerationID uuid.UUID,
) (models.Member, error) {
	initiator, err := s.repo.GetMemberByAccountAndAgglomeration(ctx, initiatorAccountID, agglomerationID)
	if err != nil {
		return models.Member{}, err
	}
	if initiator.IsNil() {
		return models.Member{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member with account id %s and agglomeration id %s not found: %w",
				initiatorAccountID, agglomerationID, err),
		)
	}

	return initiator, nil
}

func (s Service) CheckPermissionsToManageRole(
	ctx context.Context,
	accountID, agglomerationID uuid.UUID,
	rank uint,
) error {
	hasPermission, err := s.repo.CheckAccountHavePermissionByCode(
		ctx, accountID, agglomerationID,
		models.RolePermissionManageRoles.String(),
	)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s does not have permission %s", accountID, models.RolePermissionManageRoles.String()),
		)
	}

	maxRole, err := s.repo.GetAccountMaxRoleInAgglomeration(ctx, accountID, agglomerationID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get max role for member %s in agglomeration %s: %w",
				accountID, agglomerationID, err),
		)
	}

	if maxRole.Rank >= rank {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s with max role rank %d cannot manage role with rank %d",
				accountID, maxRole.Rank, rank),
		)
	}

	return nil
}

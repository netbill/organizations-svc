package role

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
	messenger messenger
}

type repo interface {
	CreateRole(ctx context.Context, params CreateParams) (entity.Role, error)

	UpdateRole(ctx context.Context, roleID uuid.UUID, params UpdateParams) (entity.Role, error)
	UpdateRoleRank(ctx context.Context, roleID uuid.UUID, newRank uint) (entity.Role, error)
	UpdateRolesRanks(
		ctx context.Context,
		agglomerationID uuid.UUID,
		order map[uint]uuid.UUID,
	) error

	DeleteRole(ctx context.Context, roleID uuid.UUID) error

	FilterRoles(
		ctx context.Context,
		filter FilterParams,
		offset uint,
		limit uint,
	) (pagi.Page[[]entity.Role], error)
	GetRole(ctx context.Context, roleID uuid.UUID) (entity.Role, error)

	GetMember(ctx context.Context, memberID uuid.UUID) (entity.Member, error)
	GetMemberByAccountAndAgglomeration(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
	) (entity.Member, error)

	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]entity.Permission, error)

	GetAccountMaxRoleInAgglomeration(ctx context.Context, accountID, agglomerationID uuid.UUID) (entity.Role, error)
	GetMemberMaxRole(ctx context.Context, memberID uuid.UUID) (entity.Role, error)

	CheckAccountHavePermissionByCode(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
		permissionKey string,
	) (bool, error)

	//CheckMemberHavePermissionByCode(ctx context.Context, memberID uuid.UUID, permissionKey string) (bool, error)
	//CheckMemberHavePermissionByID(ctx context.Context, memberID, permissionID uuid.UUID) (bool, error)
	//CheckAccountHavePermissionByCode(ctx context.Context, accountID, agglomerationID uuid.UUID, permissionKey string) (bool, error)
	//CheckAccountHavePermissionByID(ctx context.Context, accountID, agglomerationID, permissionID uuid.UUID) (bool, error)
}

type messenger interface {
}

func (s Service) GetInitiator(ctx context.Context, memberID uuid.UUID) (entity.Member, error) {
	row, err := s.repo.GetMember(ctx, memberID)
	if err != nil {
		return entity.Member{}, err
	}
	if row.IsNil() {
		return entity.Member{}, errx.ErrorMemberNotFound.Raise(
			fmt.Errorf("member with id %s not found", memberID),
		)
	}

	return row, nil
}

func (s Service) GetInitiatorByAccountAndAgglomeration(
	ctx context.Context,
	initiatorAccountID, agglomerationID uuid.UUID,
) (entity.Member, error) {
	initiator, err := s.repo.GetMemberByAccountAndAgglomeration(ctx, initiatorAccountID, agglomerationID)
	if err != nil {
		return entity.Member{}, err
	}
	if initiator.IsNil() {
		return entity.Member{}, errx.ErrorNotEnoughRights.Raise(
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
		entity.RolePermissionManageRoles.String(),
	)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s does not have permission %s", accountID, entity.RolePermissionManageRoles.String()),
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

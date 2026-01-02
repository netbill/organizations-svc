package member

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
	CreateMember(ctx context.Context, accountID, agglomerationID uuid.UUID) (models.Member, error)

	UpdateMember(ctx context.Context, ID uuid.UUID, params UpdateParams) (models.Member, error)

	GetMember(ctx context.Context, memberID uuid.UUID) (models.Member, error)
	GetMemberByAccountAndAgglomeration(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
	) (models.Member, error)
	GetMembers(
		ctx context.Context,
		filter FilterParams,
		offset uint,
		limit uint,
	) (pagi.Page[[]models.Member], error)

	DeleteMember(ctx context.Context, memberID uuid.UUID) error

	CheckMemberHavePermission(
		ctx context.Context,
		memberID uuid.UUID,
		permissionCode string,
	) (bool, error)
	GetMemberMaxRole(ctx context.Context, memberID uuid.UUID) (models.Role, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteMemberCreated(ctx context.Context, member models.Member) error
	WriteMemberUpdated(ctx context.Context, member models.Member) error
	WriteMemberDeleted(ctx context.Context, member models.Member) error
}

func (s Service) CheckAccessToManageOtherMember(
	ctx context.Context,
	firstMemberID, secMemberID uuid.UUID,
) error {
	hasPermission, err := s.repo.CheckMemberHavePermission(
		ctx,
		firstMemberID,
		models.RolePermissionManageMembers,
	)
	if err != nil {
		return err
	}

	if !hasPermission {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member %s has no manage members permission", firstMemberID),
		)
	}

	firstMaxRole, err := s.repo.GetMemberMaxRole(ctx, firstMemberID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get max role for member %s: %w", firstMemberID, err),
		)
	}
	if firstMaxRole.IsNil() {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s has no roles assigned", firstMemberID),
		)
	}

	secMaxRole, err := s.repo.GetMemberMaxRole(ctx, secMemberID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get max role for member %s: %w", secMemberID, err),
		)
	}
	if secMaxRole.IsNil() {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s has no roles assigned", secMemberID),
		)
	}

	if firstMaxRole.Rank >= secMaxRole.Rank {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf(
				"member %s with rank %d cannot manage member %s with rank %d",
				firstMemberID,
				firstMaxRole.Rank,
				secMemberID,
				secMaxRole.Rank,
			),
		)
	}

	return nil
}

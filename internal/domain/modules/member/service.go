package member

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
	CreateMember(ctx context.Context, accountID, agglomerationID uuid.UUID) (entity.Member, error)

	UpdateMember(ctx context.Context, ID uuid.UUID, params UpdateParams) (entity.Member, error)

	GetMember(ctx context.Context, memberID uuid.UUID) (entity.Member, error)
	GetMemberByAccountAndAgglomeration(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
	) (entity.Member, error)

	FilterMembers(
		ctx context.Context,
		filter FilterParams,
		offset uint,
		limit uint,
	) (pagi.Page[[]entity.Member], error)

	DeleteMember(ctx context.Context, memberID uuid.UUID) error

	CheckMemberHavePermission(
		ctx context.Context,
		memberID uuid.UUID,
		permissionCode string,
	) (bool, error)

	CanInteract(ctx context.Context, firstMemberID, secondMemberID uuid.UUID) (bool, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteMemberCreated(ctx context.Context, member entity.Member) error
	WriteMemberUpdated(ctx context.Context, member entity.Member) error
	WriteMemberDeleted(ctx context.Context, memberID uuid.UUID) error
}

func (s Service) CheckAccessToManageOtherMember(
	ctx context.Context,
	initiatorID, memberID uuid.UUID,
) error {
	hasPermission, err := s.repo.CheckMemberHavePermission(
		ctx,
		initiatorID,
		entity.RolePermissionManageMembers.String(),
	)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member %s has no manage members permission", initiatorID),
		)
	}

	canInteract, err := s.repo.CanInteract(ctx, initiatorID, memberID)
	if err != nil {
		return err
	}
	if !canInteract {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("members %s and %s cannot interact", initiatorID, memberID),
		)
	}

	return nil
}

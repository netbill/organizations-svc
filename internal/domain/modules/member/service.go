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

	FilterMembers(
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

	CanInteract(ctx context.Context, firstMemberID, secondMemberID uuid.UUID) (bool, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteMemberCreated(ctx context.Context, member models.Member) error
	WriteMemberUpdated(ctx context.Context, member models.Member) error
	WriteMemberDeleted(ctx context.Context, member models.Member) error
}

func (s Service) CheckAccessToManageOtherMember(
	ctx context.Context,
	initiatorID, memberID uuid.UUID,
) error {
	hasPermission, err := s.repo.CheckMemberHavePermission(
		ctx,
		initiatorID,
		models.RolePermissionManageMembers.String(),
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

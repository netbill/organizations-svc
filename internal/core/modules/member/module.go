package member

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/orgperm"
	"github.com/netbill/restkit/pagi"
)

type Module struct {
	repo      repo
	messenger messenger
}

func New(repo repo, messenger messenger) *Module {
	return &Module{
		repo:      repo,
		messenger: messenger,
	}
}

type repo interface {
	CreateMember(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error)
	UpdateMember(ctx context.Context, ID uuid.UUID, params UpdateParams) (models.Member, error)
	GetMember(ctx context.Context, memberID uuid.UUID) (models.Member, error)
	GetMemberByAccountAndOrganization(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (models.Member, error)
	GetMembers(
		ctx context.Context,
		filter FilterParams,
		limit uint,
		offset uint,
	) (pagi.Page[[]models.Member], error)
	DeleteMember(ctx context.Context, memberID uuid.UUID) error

	GetRole(ctx context.Context, roleID uuid.UUID) (models.Role, error)

	CheckMemberHavePermission(
		ctx context.Context,
		memberID uuid.UUID,
		permissionID uuid.UUID,
	) (bool, error)
	GetMemberMaxRoleRank(ctx context.Context, memberID uuid.UUID) (int32, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteOrgMemberCreated(ctx context.Context, member models.Member) error
	WriteOrgMemberUpdated(ctx context.Context, member models.Member) error
	WriteOrgMemberDeleted(ctx context.Context, memberID uuid.UUID) error
}

func (m *Module) checkAbilityToUpdateMember(
	ctx context.Context,
	initiator models.AccountActor,
	organizationID uuid.UUID,
	memberID uuid.UUID,
) error {
	member, err := m.GetInitiator(ctx, initiator, organizationID)
	if err != nil {
		return err
	}
	if member.Head {
		return nil
	}

	hasPermission, err := m.repo.CheckMemberHavePermission(ctx, member.AccountID, orgperm.MembersUpdateID)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member %s has no update members permission", member.ID),
		)
	}

	firstMaxRank, err := m.repo.GetMemberMaxRoleRank(ctx, member.ID)
	if err != nil {
		return fmt.Errorf("failed to get max role for member %s: %w", member.AccountID, err)
	}

	secondMaxRank, err := m.repo.GetMemberMaxRoleRank(ctx, memberID)
	if err != nil {
		return fmt.Errorf("failed to get max role for member %s: %w", memberID, err)
	}

	if firstMaxRank < secondMaxRank {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf(
				"initiator member %s has max role rank %d less than max role rank %d of member %s",
				member.AccountID, firstMaxRank, secondMaxRank, memberID,
			),
		)
	}

	return nil
}

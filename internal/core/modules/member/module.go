package member

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
	"github.com/netbill/organizations-svc/internal/core/errx"
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
	CreateMember(ctx context.Context, accountID, organizationID uuid.UUID) (domain.Member, error)
	UpdateMember(ctx context.Context, ID uuid.UUID, params UpdateParams) (domain.Member, error)
	GetMember(ctx context.Context, memberID uuid.UUID) (domain.Member, error)
	GetMemberByAccountAndOrganization(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (domain.Member, error)
	GetMembers(
		ctx context.Context,
		filter FilterParams,
		limit uint,
		offset uint,
	) (pagi.Page[[]domain.Member], error)
	DeleteMember(ctx context.Context, memberID uuid.UUID) error

	GetRole(ctx context.Context, roleID uuid.UUID) (domain.Role, error)

	CheckMemberHavePermission(
		ctx context.Context,
		memberID uuid.UUID,
		permissionID uuid.UUID,
	) (bool, error)
	GetMemberMaxRole(ctx context.Context, memberID uuid.UUID) (domain.Role, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteOrgMemberCreated(ctx context.Context, member domain.Member) error
	WriteOrgMemberUpdated(ctx context.Context, member domain.Member) error
	WriteOrgMemberDeleted(ctx context.Context, memberID uuid.UUID) error
}

func (m *Module) checkAbilityToUpdateMember(
	ctx context.Context,
	initiator domain.AccountActor,
	organizationID uuid.UUID,
	memberID uuid.UUID,
) error {
	member, err := m.GetInitiator(ctx, initiator, organizationID)
	if err != nil {
		return err
	}

	if !member.Head {
		hasPermission, err := m.repo.CheckMemberHavePermission(ctx, member.AccountID, orgperm.MembersUpdateID)
		if err != nil {
			return err
		}
		if !hasPermission {
			return errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("initiator member %s has no update members permission", member.ID),
			)
		}

		firstMaxRole, err := m.repo.GetMemberMaxRole(ctx, member.ID)
		if err != nil {
			return fmt.Errorf("failed to get max role for member %s: %w", member.AccountID, err)
		}

		secMaxRole, err := m.repo.GetMemberMaxRole(ctx, memberID)
		if err != nil {
			return fmt.Errorf("failed to get max role for member %s: %w", memberID, err)
		}

		if firstMaxRole.Rank < secMaxRole.Rank {
			return errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf(
					"member %s with rank %d cannot manage member %s with rank %d",
					member.AccountID,
					firstMaxRole.Rank,
					memberID,
					secMaxRole.Rank,
				),
			)
		}
	}

	return nil
}

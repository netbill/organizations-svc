package member

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
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

	CheckMemberHavePermission(
		ctx context.Context,
		memberID uuid.UUID,
		permissionCode string,
	) (bool, error)
	GetMemberMaxRole(ctx context.Context, memberID uuid.UUID) (models.Role, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteOrgMemberCreated(ctx context.Context, member models.Member) error
	WriteOrgMemberUpdated(ctx context.Context, member models.Member) error
	WriteOrgMemberDeleted(ctx context.Context, memberID uuid.UUID) error
}

func (m *Module) checkPermissionToInteractWithMember(
	ctx context.Context,
	initiator models.InitiatorData,
	organizationID uuid.UUID,
	memberID uuid.UUID,
) error {
	member, err := m.GetInitiatorMember(ctx, initiator, organizationID)
	if err != nil {
		return err
	}

	if !member.Head {
		hasPermission, err := m.repo.CheckMemberHavePermission(
			ctx,
			member.AccountID,
			models.RolePermissionManageMembers,
		)
		if err != nil {
			return err
		}
		if !hasPermission {
			return errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("initiator member %s has no manage members permission", member.ID),
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

package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/orgperm"
)

type CreateParams struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	Rank           int32     `json:"rank"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Color          string    `json:"color"`
}

func (m *Module) Create(
	ctx context.Context,
	initiator models.AccountActor,
	params CreateParams,
) (role models.Role, err error) {
	member, err := m.getInitiator(ctx, initiator, params.OrganizationID)
	if err != nil {
		return models.Role{}, err
	}
	if !member.Head {
		hasPermission, err := m.repo.CheckMemberHavePermission(ctx, member.ID, orgperm.RolesManageID)
		if err != nil {
			return models.Role{}, err
		}
		if !hasPermission {
			return models.Role{}, errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("member %s does not have permission %s", member.ID, orgperm.RolesManageID),
			)
		}

		initiatorMaxRank, err := m.repo.GetMemberMaxRoleRank(ctx, member.ID)
		if err != nil {
			return models.Role{}, fmt.Errorf("failed to get max role for member %s: %w", member.AccountID, err)
		}

		if initiatorMaxRank <= params.Rank {
			return models.Role{}, errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("member %s with max role rank %d cannot manage role with rank %d",
					member.ID, initiatorMaxRank, params.Rank,
				),
			)
		}
	}

	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		err = m.repo.LockOrgRoleRankRevision(ctx, params.OrganizationID)
		if err != nil {
			return err
		}

		role, err = m.repo.CreateRole(ctx, params)
		if err != nil {
			return err
		}

		ranks, err := m.repo.CreateRoleRank(ctx, role.ID, params.Rank)
		if err != nil {
			return err
		}

		if err = m.messenger.WriteOrgRoleCreated(ctx, role); err != nil {
			return err
		}

		if err = m.repo.CreateRolePermissionsRevision(ctx, role.ID); err != nil {
			return err
		}

		revision, err := m.repo.BumpOrgRoleRankRevision(ctx, params.OrganizationID)
		if err != nil {
			return err
		}

		if err = m.messenger.WriteOrgRolesRanksUpdated(ctx, params.OrganizationID, ranks, revision); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Role{}, err
	}

	return role, nil
}

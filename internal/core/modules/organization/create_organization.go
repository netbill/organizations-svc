package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type CreateParams struct {
	Name string
}

func (s Service) CreateOrganization(
	ctx context.Context,
	accountID uuid.UUID,
	params CreateParams,
) (org models.Organization, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		org, err = s.repo.CreateOrganization(ctx, params)
		if err != nil {
			return fmt.Errorf("failed to create organization, cause: %w", err)
		}

		err = s.messenger.WriteOrganizationCreated(ctx, org)
		if err != nil {
			return fmt.Errorf("failed to publish organization create event, cause: %w", err)
		}

		role, err := s.createRoleHead(ctx, org.ID)
		if err != nil {
			return err
		}

		if _, err = s.createMemberHead(ctx, accountID, org.ID, role.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Organization{}, err
	}

	return org, err
}

func (s Service) createRoleHead(
	ctx context.Context,
	organizationID uuid.UUID,
) (role models.Role, err error) {
	role, err = s.repo.CreateHeadRole(ctx, organizationID)
	if err != nil {
		return models.Role{}, fmt.Errorf("failed to create role, cause: %w", err)
	}

	err = s.messenger.WriteOrgRoleCreated(ctx, role)
	if err != nil {
		return models.Role{}, fmt.Errorf("failed to publish role create event, cause: %w", err)
	}

	per, err := s.repo.GetRolePermissions(ctx, role.ID)
	if err != nil {
		return models.Role{}, fmt.Errorf("failed to get role permissions, cause: %w", err)
	}

	if err = s.messenger.WriteOrgRolePermissionsUpdated(ctx, role, per); err != nil {
		return models.Role{}, fmt.Errorf("failed to publish role permissions updated event, cause: %w", err)
	}

	return role, nil
}

func (s Service) createMemberHead(
	ctx context.Context,
	accountID uuid.UUID,
	organizationID uuid.UUID,
	roleID uuid.UUID,
) (member models.Member, err error) {
	member, err = s.repo.CreateMember(ctx, accountID, organizationID)
	if err != nil {
		return models.Member{}, fmt.Errorf("failed to create member, cause: %w", err)
	}

	err = s.repo.AddMemberRole(ctx, member.ID, roleID)
	if err != nil {
		return models.Member{}, fmt.Errorf("failed to assign head role to member, cause: %w", err)
	}

	err = s.messenger.WriteOrgMemberCreated(ctx, member)
	if err != nil {
		return models.Member{}, fmt.Errorf("failed to publish member create event, cause: %w", err)
	}

	return member, nil
}

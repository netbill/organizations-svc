package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type CreateParams struct {
	Name string
	Icon *string
}

func (s Service) CreateOrganization(
	ctx context.Context,
	accountID uuid.UUID,
	params CreateParams,
) (org models.Organization, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		org, err = s.repo.CreateOrganization(ctx, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create organization: %w", err),
			)
		}

		err = s.messenger.WriteOrganizationCreated(ctx, org)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish organization create event: %w", err),
			)
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

func (s Service) createRoleHead(ctx context.Context, organizationID uuid.UUID) (role models.Role, err error) {
	role, err = s.repo.CreateHeadRole(ctx, organizationID)
	if err != nil {
		return models.Role{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create role: %w", err),
		)
	}

	err = s.messenger.WriteRoleCreated(ctx, role)
	if err != nil {
		return models.Role{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish role create event: %w", err),
		)
	}

	per, err := s.repo.GetRolePermissions(ctx, role.ID)
	if err != nil {
		return models.Role{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get role permissions: %w", err),
		)
	}

	if err = s.messenger.WriteRolePermissionsUpdated(
		ctx,
		role.ID,
		per,
	); err != nil {
		return models.Role{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish role permissions updated event: %w", err),
		)
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
		return models.Member{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create member: %w", err),
		)
	}

	err = s.repo.AddMemberRole(ctx, member.ID, roleID)
	if err != nil {
		return models.Member{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to assign head role to member: %w", err),
		)
	}

	err = s.messenger.WriteMemberCreated(ctx, member)
	if err != nil {
		return models.Member{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish member create event: %w", err),
		)
	}

	return member, nil
}

package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/domain/errx"
	"github.com/netbill/organizations-svc/internal/domain/models"
)

type CreateParams struct {
	Name string
	Icon *string
}

func (s Service) CreateOrganization(
	ctx context.Context,
	accountID uuid.UUID,
	params CreateParams,
) (agglo models.Organization, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		agglo, err = s.repo.CreateOrganization(ctx, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create organization: %w", err),
			)
		}

		err = s.messenger.WriteOrganizationCreated(ctx, agglo)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish organization create event: %w", err),
			)
		}

		role, err := s.createRoleHead(ctx, agglo.ID)
		if err != nil {
			return err
		}

		_, err = s.createMemberHead(ctx, accountID, agglo.ID, role.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Organization{}, err
	}

	return agglo, err
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

	return member, nil
}

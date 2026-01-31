package organization

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
)

func (s Service) DeleteOrganization(ctx context.Context, accountID, organizationID uuid.UUID) error {
	organization, err := s.GetOrganization(ctx, organizationID)
	if err != nil {
		return err
	}

	initiator, err := s.repo.GetMemberByAccountAndOrganization(ctx, accountID, organization.ID)
	if err != nil {
		return err
	}

	role, err := s.repo.GetMemberMaxRole(ctx, initiator.ID)
	if err != nil {
		if errors.Is(err, errx.ErrorRoleNotFound) {
			return errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("member with id %s has no role in organization %s: %w",
					initiator.AccountID, organization.ID, err.Error()),
			)
		}
		return err
	}

	if role.Head != true {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("only organization head can delete organization"),
		)
	}

	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		err = s.repo.DeleteOrganization(ctx, organizationID)
		if err != nil {
			return fmt.Errorf("failed to delete organization: %w", err)
		}

		err = s.messenger.WriteOrganizationDeleted(ctx, organization)
		if err != nil {
			return fmt.Errorf("failed to publish organization delete event: %w", err)
		}

		return nil
	})
}

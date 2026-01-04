package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/pagi"
)

func (s Service) GetOrganization(ctx context.Context, organizationID uuid.UUID) (models.Organization, error) {
	res, err := s.repo.GetOrganizationByID(ctx, organizationID)
	if err != nil {
		return models.Organization{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get organization by id: %w", err),
		)
	}
	if res.IsNil() {
		return models.Organization{}, errx.ErrorOrganizationNotFound.Raise(
			fmt.Errorf("organization with id %s not found", organizationID),
		)
	}

	return res, nil
}

type FilterParams struct {
	Name   *string `json:"name,omitempty"`
	Status *string `json:"status,omitempty"`
}

func (s Service) GetOrganizations(
	ctx context.Context,
	params FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Organization], error) {
	res, err := s.repo.GetOrganizations(ctx, params, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Organization]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("filter organizations: %w", err),
		)
	}

	return res, nil
}

func (s Service) GetOrganizationForUser(
	ctx context.Context,
	accountID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Organization], error) {
	res, err := s.repo.GetOrganizationsForUser(ctx, accountID, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Organization]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("get organization for user: %w", err),
		)
	}

	return res, nil
}

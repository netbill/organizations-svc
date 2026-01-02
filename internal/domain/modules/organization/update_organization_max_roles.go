package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/domain/errx"
	"github.com/netbill/organizations-svc/internal/domain/models"
)

func (s Service) UpdateOrganizationMaxRoles(
	ctx context.Context,
	organizationID uuid.UUID,
	maxRoles uint,
) (models.Organization, error) {
	agglo, err := s.repo.UpdateOrganizationMaxRoles(ctx, organizationID, maxRoles)
	if err != nil {
		return models.Organization{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update organization max roles: %w", err),
		)
	}

	return agglo, nil
}

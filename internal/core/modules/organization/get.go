package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

func (m *Module) GetByID(
	ctx context.Context,
	organizationID uuid.UUID,
) (models.Organization, error) {
	res, err := m.repo.GetOrganizationByID(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	return res, nil
}

type FilterParams struct {
	Text   *string `json:"name,omitempty"`
	Status *string `json:"status,omitempty"`
}

func (m *Module) GetList(
	ctx context.Context,
	params FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Organization], error) {
	res, err := m.repo.GetOrganizations(ctx, params, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Organization]{}, err
	}

	return res, nil
}

func (m *Module) GetForUser(
	ctx context.Context,
	actor models.AccountActor,
	limit, offset uint,
) (pagi.Page[[]models.Organization], error) {
	res, err := m.repo.GetOrganizationsForUser(ctx, actor, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Organization]{}, err
	}

	return res, nil
}

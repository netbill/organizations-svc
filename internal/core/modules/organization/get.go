package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
	"github.com/netbill/restkit/pagi"
)

func (m *Module) GetByID(
	ctx context.Context,
	organizationID uuid.UUID,
) (domain.Organization, error) {
	res, err := m.repo.GetOrganizationByID(ctx, organizationID)
	if err != nil {
		return domain.Organization{}, err
	}

	return res, nil
}

type FilterParams struct {
	Name   *string `json:"name,omitempty"`
	Status *string `json:"status,omitempty"`
}

func (m *Module) GetList(
	ctx context.Context,
	params FilterParams,
	limit, offset uint,
) (pagi.Page[[]domain.Organization], error) {
	res, err := m.repo.GetOrganizations(ctx, params, limit, offset)
	if err != nil {
		return pagi.Page[[]domain.Organization]{}, err
	}

	return res, nil
}

func (m *Module) GetForUser(
	ctx context.Context,
	initiator domain.AccountActor,
	limit, offset uint,
) (pagi.Page[[]domain.Organization], error) {
	res, err := m.repo.GetOrganizationsForUser(ctx, initiator, limit, offset)
	if err != nil {
		return pagi.Page[[]domain.Organization]{}, err
	}

	return res, nil
}

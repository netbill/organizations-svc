package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
	"github.com/netbill/restkit/pagi"
)

type FilterParams struct {
	OrganizationID *uuid.UUID
	RolesID        *[]uuid.UUID
	Rank           *int
	Name           *string
}

func (m *Module) GetList(
	ctx context.Context,
	params FilterParams,
	limit, offset uint,
) (pagi.Page[[]domain.Role], error) {
	res, err := m.repo.GetRoles(ctx, params, limit, offset)
	if err != nil {
		return pagi.Page[[]domain.Role]{}, err
	}

	return res, nil
}

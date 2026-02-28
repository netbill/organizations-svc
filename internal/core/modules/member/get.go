package member

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

func (m *Module) GetByID(
	ctx context.Context,
	memberID uuid.UUID,
) (models.Member, error) {
	return m.repo.GetMember(ctx, memberID)
}

func (m *Module) GetByAccountAndOrgs(
	ctx context.Context,
	actor models.AccountActor,
	organizationIDs []uuid.UUID,
) ([]models.Member, error) {
	rows, err := m.repo.GetMembersByAccountAndOrgs(ctx, actor, organizationIDs)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (m *Module) GetByAccountAndOrganization(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.Member, error) {
	row, err := m.repo.GetMemberByAccountAndOrganization(ctx, actor, organizationID)
	if err != nil {
		return models.Member{}, err
	}

	return row, nil
}

type FilterParams struct {
	OrganizationID *uuid.UUID
	AccountID      *uuid.UUID
	Head           *bool
	Username       *string
	BestMatch      *string
	Label          *string
	Position       *string
}

func (m *Module) GetList(
	ctx context.Context,
	filter FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Member], error) {
	res, err := m.repo.GetMembers(ctx, filter, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Member]{}, err
	}

	return res, nil
}

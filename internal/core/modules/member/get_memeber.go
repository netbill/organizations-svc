package member

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

func (m *Module) GetMemberByID(ctx context.Context, memberID uuid.UUID) (models.Member, error) {
	row, err := m.repo.GetMember(ctx, memberID)
	if err != nil {
		return models.Member{}, err
	}

	return row, nil
}

func (m *Module) GetMemberByAccountAndOrganization(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (models.Member, error) {
	row, err := m.repo.GetMemberByAccountAndOrganization(ctx, accountID, organizationID)
	if err != nil {
		return models.Member{}, err
	}

	return row, nil
}

func (m *Module) GetInitiatorMember(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (models.Member, error) {
	initiator, err := m.GetMemberByAccountAndOrganization(ctx, accountID, organizationID)
	if errors.Is(err, errx.ErrorMemberNotFound) {
		return models.Member{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member with account id %s and organization id %s not found: %w",
				accountID, organizationID, err.Error()),
		)
	}

	return initiator, err
}

type FilterParams struct {
	OrganizationID *uuid.UUID
	AccountID      *uuid.UUID
	RoleID         *uuid.UUID
	Username       *string
	BestMatch      *string
	PermissionCode *string
	Label          *string
	Position       *string
	RoleRankUp     *uint
	RoleRankDown   *uint
}

func (m *Module) GetMembers(
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

package member

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/restkit/pagi"
)

func (m *Module) GetByID(
	ctx context.Context,
	memberID uuid.UUID,
) (domain.Member, error) {
	row, err := m.repo.GetMember(ctx, memberID)
	if err != nil {
		return domain.Member{}, err
	}

	return row, nil
}

func (m *Module) GetInitiator(
	ctx context.Context,
	initiator domain.AccountActor,
	organizationID uuid.UUID,
) (domain.Member, error) {
	member, err := m.GetByAccountAndOrganization(ctx, initiator, organizationID)
	if errors.Is(err, errx.ErrorMemberNotFound) {
		return domain.Member{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf(
				"initiator member with account id %s and organization id %s not found: %w",
				member.AccountID, organizationID, err,
			),
		)
	}

	return member, err
}

func (m *Module) GetByAccountAndOrganization(
	ctx context.Context,
	initiator domain.AccountActor,
	organizationID uuid.UUID,
) (domain.Member, error) {
	row, err := m.repo.GetMemberByAccountAndOrganization(ctx, initiator, organizationID)
	if err != nil {
		return domain.Member{}, err
	}

	return row, nil
}

type FilterParams struct {
	OrganizationID *uuid.UUID
	AccountID      *uuid.UUID
	RoleID         *uuid.UUID
	Head           *bool
	Username       *string
	BestMatch      *string
	PermissionID   *uuid.UUID
	Label          *string
	Position       *string
	RoleRankUp     *uint
	RoleRankDown   *uint
}

func (m *Module) GetList(
	ctx context.Context,
	filter FilterParams,
	limit, offset uint,
) (pagi.Page[[]domain.Member], error) {
	res, err := m.repo.GetMembers(ctx, filter, limit, offset)
	if err != nil {
		return pagi.Page[[]domain.Member]{}, err
	}

	return res, nil
}

package member

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/domain/errx"
	"github.com/netbill/organizations-svc/internal/domain/models"
	"github.com/netbill/pagi"
)

func (s Service) GetMemberByID(ctx context.Context, memberID uuid.UUID) (models.Member, error) {
	row, err := s.repo.GetMember(ctx, memberID)
	if err != nil {
		return models.Member{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get member with id %s: %w", memberID, err),
		)
	}
	if row.IsNil() {
		return models.Member{}, errx.ErrorMemberNotFound.Raise(
			fmt.Errorf("member with id %s not found", memberID),
		)
	}

	return row, nil
}

func (s Service) GetMemberByAccountAndOrganization(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (models.Member, error) {
	row, err := s.repo.GetMemberByAccountAndOrganization(ctx, accountID, organizationID)
	if err != nil {
		return models.Member{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get member with account id %s and organization id %s: %w",
				accountID, organizationID, err),
		)
	}
	if row.IsNil() {
		return models.Member{}, errx.ErrorMemberNotFound.Raise(
			fmt.Errorf("member with account id %s and organization id %s not found", accountID, organizationID),
		)
	}

	return row, nil
}

func (s Service) GetInitiatorMember(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (models.Member, error) {
	initiator, err := s.GetMemberByAccountAndOrganization(ctx, accountID, organizationID)
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

func (s Service) GetMembers(
	ctx context.Context,
	filter FilterParams,
	offset uint,
	limit uint,
) (pagi.Page[[]models.Member], error) {
	res, err := s.repo.GetMembers(ctx, filter, offset, limit)
	if err != nil {
		return pagi.Page[[]models.Member]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to filter members: %w", err),
		)
	}

	return res, nil
}

package member

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/pagi"
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

func (s Service) GetMemberByAccountAndAgglomeration(
	ctx context.Context,
	accountID, agglomerationID uuid.UUID,
) (models.Member, error) {
	row, err := s.repo.GetMemberByAccountAndAgglomeration(ctx, accountID, agglomerationID)
	if err != nil {
		return models.Member{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get member with account id %s and agglomeration id %s: %w",
				accountID, agglomerationID, err),
		)
	}
	if row.IsNil() {
		return models.Member{}, errx.ErrorMemberNotFound.Raise(
			fmt.Errorf("member with account id %s and agglomeration id %s not found", accountID, agglomerationID),
		)
	}

	return row, nil
}

func (s Service) GetInitiatorMember(
	ctx context.Context,
	accountID, agglomerationID uuid.UUID,
) (models.Member, error) {
	initiator, err := s.GetMemberByAccountAndAgglomeration(ctx, accountID, agglomerationID)
	if errors.Is(err, errx.ErrorMemberNotFound) {
		return models.Member{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member with account id %s and agglomeration id %s not found: %w",
				accountID, agglomerationID, err.Error()),
		)
	}

	return initiator, err
}

type FilterParams struct {
	AgglomerationID *uuid.UUID
	AccountID       *uuid.UUID
	RoleID          *uuid.UUID
	Username        *string
	BestMatch       *string
	PermissionCode  *string
	Label           *string
	Position        *string
	RoleRankUp      *uint
	RoleRankDown    *uint
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

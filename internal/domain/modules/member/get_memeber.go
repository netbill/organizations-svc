package member

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

func (s Service) GetMember(ctx context.Context, memberID uuid.UUID) (entity.Member, error) {
	row, err := s.repo.GetMember(ctx, memberID)
	if err != nil {
		return entity.Member{}, err
	}
	if row.IsNil() {
		return entity.Member{}, errx.ErrorMemberNotFound.Raise(
			fmt.Errorf("member with id %s not found", memberID),
		)
	}

	return row, nil
}

func (s Service) GetMemberByAccountAndAgglomeration(
	ctx context.Context,
	accountID, agglomerationID uuid.UUID,
) (entity.Member, error) {
	row, err := s.repo.GetMemberByAccountAndAgglomeration(ctx, accountID, agglomerationID)
	if err != nil {
		return entity.Member{}, err
	}
	if row.IsNil() {
		return entity.Member{}, errx.ErrorMemberNotFound.Raise(
			fmt.Errorf("member with account id %s and agglomeration id %s not found", accountID, agglomerationID),
		)
	}

	return row, nil
}

func (s Service) GetInitiatorMemberByAccountAndAgglomeration(
	ctx context.Context,
	accountID, agglomerationID uuid.UUID,
) (entity.Member, error) {
	initiator, err := s.GetMemberByAccountAndAgglomeration(ctx, accountID, agglomerationID)
	if errors.Is(err, errx.ErrorMemberNotFound) {
		return entity.Member{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member with account id %s and agglomeration id %s not found: %w",
				accountID, agglomerationID, err.Error()),
		)
	}

	return initiator, err
}

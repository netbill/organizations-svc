package invite

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/pagi"
)

func (s Service) getInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error) {
	res, err := s.repo.GetInvite(ctx, ID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get invite by ID %s: %w", ID.String(), err),
		)
	}
	if res.IsNil() {
		return models.Invite{}, errx.ErrorInviteNotFound.Raise(
			fmt.Errorf("invite with ID %s not found", ID.String()),
		)
	}

	return res, nil
}

func (s Service) GetInvite(
	ctx context.Context,
	accountID, ID uuid.UUID,
) (models.Invite, error) {
	res, err := s.getInvite(ctx, ID)
	if err != nil {
		return models.Invite{}, err
	}

	if res.AccountID != accountID {
		member, err := s.repo.GetMemberByAccountAndAgglomeration(
			ctx,
			accountID,
			res.AgglomerationID,
		)
		if err != nil {
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get member by account %s and agglomeration %s: %w",
					accountID.String(), res.AgglomerationID.String(), err),
			)
		}
		if member.IsNil() {
			return models.Invite{}, errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("account is not a member of the agglomeration"),
			)
		}
	}

	return res, nil
}

func (s Service) GetAgglomerationInvites(
	ctx context.Context,
	accountID, agglomerationID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	_, err := s.getInitiator(ctx, accountID, agglomerationID)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, err
	}

	res, err := s.repo.GetAgglomerationInvites(ctx, agglomerationID, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get agglomeration invites: %w", err),
		)
	}

	return res, nil
}

func (s Service) GetAccountInvites(
	ctx context.Context,
	accountID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	res, err := s.repo.GetAccountInvites(ctx, accountID, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get account invites: %w", err),
		)
	}

	return res, nil
}

package invite

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/domain/errx"
	"github.com/netbill/organizations-svc/internal/domain/models"
	"github.com/netbill/pagi"
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
		member, err := s.repo.GetMemberByAccountAndOrganization(
			ctx,
			accountID,
			res.OrganizationID,
		)
		if err != nil {
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get member by account %s and organization %s: %w",
					accountID.String(), res.OrganizationID.String(), err),
			)
		}
		if member.IsNil() {
			return models.Invite{}, errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("account is not a member of the organization"),
			)
		}
	}

	return res, nil
}

func (s Service) GetOrganizationInvites(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	_, err := s.getInitiator(ctx, accountID, organizationID)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, err
	}

	res, err := s.repo.GetOrganizationInvites(ctx, organizationID, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get organization invites: %w", err),
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

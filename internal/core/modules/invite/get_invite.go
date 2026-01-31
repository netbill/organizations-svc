package invite

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

func (s Service) GetInviteForAccount(
	ctx context.Context,
	accountID, ID uuid.UUID,
) (models.Invite, error) {
	res, err := s.repo.GetInvite(ctx, ID)
	if err != nil {
		return models.Invite{}, err
	}

	if res.AccountID != accountID {
		_, err = s.repo.GetMemberByAccountAndOrganization(
			ctx,
			accountID,
			res.OrganizationID,
		)
		if err != nil {
			if errors.Is(err, errx.ErrorMemberNotFound) {
				return models.Invite{}, errx.ErrorInviteNotFound.Raise(
					fmt.Errorf("account has no rights to view this invite"),
				)
			}
			return models.Invite{}, err
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
		return pagi.Page[[]models.Invite]{}, err
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
		return pagi.Page[[]models.Invite]{}, err
	}

	return res, nil
}

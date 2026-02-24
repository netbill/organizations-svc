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

func (m *Module) GetListForOrganization(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	_, err := m.getInitiator(ctx, actor, organizationID)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, err
	}

	res, err := m.repo.GetOrganizationInvites(ctx, organizationID, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, err
	}

	return res, nil
}

func (m *Module) GetListForAccount(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	_, err := m.getInitiator(ctx, actor, organizationID)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, err
	}

	res, err := m.repo.GetOrganizationInvites(ctx, organizationID, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, err
	}

	return res, nil
}

func (m *Module) GetForAccount(
	ctx context.Context,
	accountID uuid.UUID,
	inviteID uuid.UUID,
) (models.Invite, error) {
	res, err := m.repo.GetInvite(ctx, inviteID)
	if err != nil {
		return models.Invite{}, err
	}

	if res.AccountID != accountID {
		_, err = m.repo.GetMemberByAccountAndOrganization(ctx, accountID, res.OrganizationID)
		if errors.Is(err, errx.ErrorMemberNotFound) {
			return models.Invite{}, errx.ErrorInviteNotFound.Raise(
				fmt.Errorf("account has no rights to view this invite"),
			)
		} else if err != nil {
			return models.Invite{}, err
		}
	}

	return res, nil
}

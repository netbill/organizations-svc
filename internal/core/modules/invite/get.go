package invite

import (
	"context"

	"github.com/google/uuid"
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
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	res, err := m.repo.GetAccountInvites(ctx, actor, limit, offset)
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
		_, err = m.getInitiator(ctx, accountID, res.OrganizationID)
		if err != nil {
			return models.Invite{}, err
		}
	}

	return res, nil
}
